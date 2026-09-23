package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"homeopathy-platform/internal/models"
	"homeopathy-platform/pkg/response"

	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartHandler struct {
	db *gorm.DB
}

const guestCartCookieName = "guest_cart_token"

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

func (t *CartHandler) Get(c *fiber.Ctx) error {
	cart, err := t.getCart(c)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to fetch cart")
	}

	if cart.ID == uuid.Nil {
		return response.OK(c, fiber.Map{"items": []models.CartItem{}})
	}

	return response.OK(c, cart)
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
	return userId, nil
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

	guestToken := c.Cookies(guestCartCookieName)

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
		Name:     guestCartCookieName,
		Value:    token,
		Path:     "/",
		HTTPOnly: true,
		Secure:   false, // true in production
		SameSite: "Lax",
		MaxAge:   60 * 60 * 24 * 30,
	})

	return cart, nil

}

func (h *CartHandler) AddItem(c *fiber.Ctx) error {
	var req addCartItemRequest

	if err := c.BodyParser(&req); err != nil {
		fmt.Println("body is being parsed")
		return response.Error(c, 400, "Invalid request body")
	}

	if req.ProductId == uuid.Nil {
		return response.Error(c, 400, "Product_id is required")
	}

	if req.Quantity <= 0 {
		return response.Error(c, 400, "Qunatity should be greater than zero")
	}

	var product models.Product

	if err := h.db.Where("id=? AND is_active=?", req.ProductId, true).First(&product).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.Error(c, 404, "Product not found")
		}
		return response.Error(
			c,
			fiber.StatusInternalServerError,
			"failed to fetch product",
		)
	}

	if req.Quantity > product.StockQty {
		return response.Error(
			c,
			fiber.StatusBadRequest,
			"requested quantity exceeds available stock",
		)
	}
	cart, err := h.getOrCreateCart(c)
	if err != nil {
		return response.Error(c, 500, "Failed to get cart")
	}
	err = h.db.Transaction(func(tx *gorm.DB) error {
		var cartItem models.CartItem

		err := tx.Where("cart_id=? AND product_id=?", cart.ID, product.ID).First(&cartItem).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			cartItem := models.CartItem{
				CartID:    cart.ID,
				ProductID: req.ProductId,
				Quantity:  req.Quantity,
			}
			return tx.Create(&cartItem).Error
		}
		if err != nil {
			return err
		}

		newQuantity := req.Quantity + cartItem.Quantity

		if product.StockQty < newQuantity {
			return errors.New("requested quantity exceeds available stock")
		}

		return tx.Model(&cartItem).Update("quantity", newQuantity).Error
	})
	if err != nil {
		if err.Error() == "requested quantity exceeds available stock" {
			return response.Error(c, 400, err.Error())
		}
		return response.Error(
			c,
			fiber.StatusInternalServerError,
			"failed to add item to cart",
		)
	}
	cart, err = h.getCart(c)
	if err != nil {
		return response.Error(
			c,
			fiber.StatusInternalServerError,
			"failed to fetch updated cart",
		)
	}
	return response.Created(c, cart)
}

func (h *CartHandler) getCart(c *fiber.Ctx) (models.Cart, error) {

	userID, err := getCurrentUserID(c)
	if err != nil {
		return models.Cart{}, err
	}

	// ==========================================
	// LOGGED-IN USER
	// ==========================================

	if userID != uuid.Nil {
		var cart models.Cart

		err := h.db.Preload("Items.Product").Where("user_id = ?", userID).First(&cart).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Cart{}, nil
		}

		if err != nil {
			return models.Cart{}, err
		}

		return cart, nil
	}

	// ==========================================
	// GUEST USER
	// ==========================================

	guestToken := c.Cookies(guestCartCookieName)

	// No cookie means no cart.
	if guestToken == "" {
		return models.Cart{}, nil
	}

	var cart models.Cart

	err = h.db.Preload("Items.Product").Where("guest_token = ?", guestToken).First(&cart).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.Cart{}, nil
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
