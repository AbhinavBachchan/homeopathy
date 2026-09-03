package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"homeopathy-platform/internal/models"

	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartHandler struct {
	db *gorm.DB
}

type addCartItemRequest struct {
	ProductId uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
}

type updateCartItemRequest struct {
	Quantity int `json:"quantity"`
}

func NewCartHandler(db *gorm.DB) *CartHandler {
	return &CartHandler{db: db}
}

func getCurrentUserID(c *fiber.Ctx) (uuid.UUID, error) {
	userIdValue := c.Locals("user_id")

	if userIdValue == nil {
		return uuid.Nil, nil
	}

	userIdStr, ok := userIdValue.(string)
	if !ok {
		return uuid.Nil, errors.New("Invalid user")
	}

	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		return uuid.Nil, err
	}
	return userId, err
}

func (h *CartHandler) getOrCreateCart(c *fiber.Ctx) (models.Cart, error) {

	// First try to find an existing cart.
	cart, err := h.getCart(c)
	if err != nil {
		return models.Cart{}, err
	}

	// Existing cart found.
	if cart.ID != uuid.Nil {
		return cart, nil
	}

	// No cart found. Determine whether this is
	// a logged-in user or a guest.
	userID, err := getCurrentUserID(c)
	if err != nil {
		return models.Cart{}, err
	}

	// ==========================================
	// LOGGED-IN USER
	// ==========================================

	if userID != uuid.Nil {
		cart = models.Cart{
			UserID: userID,
		}

		if err := h.db.Create(&cart).Error; err != nil {
			return models.Cart{}, err
		}

		return cart, nil
	}

	// ==========================================
	// GUEST USER
	// ==========================================

	guestToken := c.Cookies("guestCartCookieName")

	// Existing cookie but no corresponding cart.
	// Reuse the token.
	if guestToken != "" {
		cart = models.Cart{
			GuestToken: guestToken,
		}

		if err := h.db.Create(&cart).Error; err != nil {
			return models.Cart{}, err
		}

		return cart, nil
	}

	// Completely new guest.
	token, err := generateGuestToken()
	if err != nil {
		return models.Cart{}, err
	}

	cart = models.Cart{
		GuestToken: token,
	}

	if err := h.db.Create(&cart).Error; err != nil {
		return models.Cart{}, err
	}

	// Save token in the browser.
	c.Cookie(&fiber.Cookie{
		Name:     "guestCartCookieName",
		Value:    token,
		Path:     "/",
		HTTPOnly: true,
		Secure:   false, // true in production
		SameSite: "Lax",
		MaxAge:   60 * 60 * 24 * 30,
	})

	return cart, nil

}

func (h *CartHandler) getCart(c *fiber.Ctx,) (models.Cart, error) {

	userID, err := getCurrentUserID(c)
	if err != nil {
		return models.Cart{}, err
	}

	// ==========================================
	// LOGGED-IN USER
	// ==========================================

	if userID != uuid.Nil {
		var cart models.Cart

		err := h.db.Where("user_id = ?", userID).First(&cart).Error

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Cart{}, err
		}

		if err != nil {
			return models.Cart{}, err
		}

		return cart, nil
	}

	// ==========================================
	// GUEST USER
	// ==========================================

	guestToken := c.Cookies("guestCartCookieName")

	// No cookie means no cart.
	if guestToken == "" {
		return models.Cart{}, nil
	}

	var cart models.Cart

	err = h.db.Where("guest_token = ?", guestToken).First(&cart).Error

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return models.Cart{}, err
	}

	if err != nil {
		return models.Cart{}, err
	}

	return cart, nil
}

func generateGuestToken() (string, error) {
	bytes := make([]byte, 32)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
