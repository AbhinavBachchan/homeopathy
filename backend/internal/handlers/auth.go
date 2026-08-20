package handlers

import (
	"time"

	"homeopathy-platform/internal/config"
	"homeopathy-platform/internal/middleware"
	"homeopathy-platform/internal/models"
	"homeopathy-platform/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewAuthHandler(db *gorm.DB, cfg *config.Config) *AuthHandler {
	return &AuthHandler{db: db, cfg: cfg}
}

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
	Phone    string `json:"phone"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.Error(c, 500, "could not process password")
		return
	}

	user := models.User{
		Email:        req.Email,
		PasswordHash: string(hash),
		Name:         req.Name,
		Phone:        req.Phone,
		Role:         models.RolePatient,
	}

	if err := h.db.Create(&user).Error; err != nil {
		response.Error(c, 409, "email already registered")
		return
	}

	token, err := h.issueToken(user)
	if err != nil {
		response.Error(c, 500, "could not issue token")
		return
	}

	response.Created(c, gin.H{"user": user, "token": token})
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	var user models.User
	if err := h.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		response.Error(c, 401, "invalid credentials")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		response.Error(c, 401, "invalid credentials")
		return
	}

	token, err := h.issueToken(user)
	if err != nil {
		response.Error(c, 500, "could not issue token")
		return
	}

	response.OK(c, gin.H{"user": user, "token": token})
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

// TODO Phase 1 P0 remaining: mobile OTP login via MSG91, Google OAuth callback.
// These slot in as additional handlers (LoginWithOTP, GoogleCallback) using the
// same issueToken() helper so they return the same JWT shape.
