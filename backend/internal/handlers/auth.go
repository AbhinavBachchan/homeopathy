package handlers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	"homeopathy-platform/internal/config"
	"homeopathy-platform/internal/middleware"
	"homeopathy-platform/internal/models"
	"homeopathy-platform/pkg/mailer"
	"homeopathy-platform/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db     *gorm.DB
	cfg    *config.Config
	mailer *mailer.Mailer
}

func NewAuthHandler(db *gorm.DB, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		db:     db,
		cfg:    cfg,
		mailer: mailer.NewMailer(cfg.BrevoAPIKey, cfg.BrevoSenderEmail, cfg.BrevoSenderName),
	}
}

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
	Phone    string `json:"phone"`
}

func (h *AuthHandler) Register(c *fiber.Ctx) error{
	var req registerRequest
	if err := c.BodyParser(&req); err != nil {
	  return c.Status(400).JSON(fiber.Map{"error":  err.Error()})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "could not process password"})
		
	}

	user := models.User{
		Email:        req.Email,
		PasswordHash: string(hash),
		Name:         req.Name,
		Phone:        req.Phone,
		Role:         models.RolePatient,
	}

	if err := h.db.Create(&user).Error; err != nil {
		return c.Status(409).JSON(fiber.Map{"error":"email already registered"})
	}

	token, err := h.issueToken(user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error":"could not issue token"})
		
	}

	return c.Status(201).JSON(fiber.Map{"user": user, "token": token})
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *fiber.Ctx) error{
	var req loginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error":err.Error()})
	}

	var user models.User
	if err := h.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{"error":"invalid credentials"})
		
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{"error":"invalid credentials"})
	}

	token, err := h.issueToken(user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error":"could not issue token"})
	}

	return c.JSON(fiber.Map{"user": user, "token": token})
}

func (h *AuthHandler) issueToken(user models.User) (string, error) {
	claims := middleware.Claims{
		UserID: user.ID.String(),
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.cfg.JWTSecret))
}

type forgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func (h *AuthHandler) ForgotPassword(c *fiber.Ctx) error{
	var req forgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, 400, err.Error())
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))

	var user models.User
	if err := h.db.Where("LOWER(email) = ?", email).First(&user).Error; err != nil {
		return response.Error(c, 404, "Email is not registered")
	}

	rawToken, err := generateResetCode()
	if err != nil {
		return response.Error(c, 500, "could not process password reset request")
	}

	hashedToken := hashResetToken(rawToken)
	expiresAt := time.Now().Add(15 * time.Minute)

	if err := h.db.Model(&user).Updates(map[string]interface{}{
		"reset_token":            hashedToken,
		"reset_token_expires_at": expiresAt,
	}).Error; err != nil {
		return response.Error(c, 500, "could not process password reset request")
	}

	// Dispatch email with plaintext reset code via Brevo mailer (or local log if unconfigured)
	go func(recipientEmail, recipientName, code string) {
		_ = h.mailer.SendPasswordResetEmail(recipientEmail, recipientName, code)
	}(user.Email, user.Name, rawToken)

	// Note: Plaintext reset token is NEVER exposed in the API response.
	return response.OK(c, fiber.Map{
		"message": "A 6 digit verification code has been sent to the registered email. Kindly enter the code.",
	})
}

type resetPasswordRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error{
	var req resetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, 400, err.Error())
		
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	token := strings.TrimSpace(req.Token)
	hashedToken := hashResetToken(token)

	var user models.User
	if err := h.db.Where("LOWER(email) = ? AND reset_token = ?", email, hashedToken).First(&user).Error; err != nil {
		return response.Error(c, 400, "invalid or expired reset code")
		
	}

	if user.ResetTokenExpiresAt == nil || time.Now().After(*user.ResetTokenExpiresAt) {
		return response.Error(c, 400, "reset code has expired, please request a new one")
		
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return response.Error(c, 500, "could not process password")
		
	}

	// Update user password and clear reset token fields
	if err := h.db.Model(&user).Updates(map[string]interface{}{
		"password_hash":          string(hash),
		"reset_token":            "",
		"reset_token_expires_at": nil,
	}).Error; err != nil {
		return response.Error(c, 500, "could not update password")
		
	}

	return response.OK(c, fiber.Map{
		"message": "Password has been successfully reset. You can now login with your new password.",
	})
}

func generateResetCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()+100000), nil
}

func hashResetToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

