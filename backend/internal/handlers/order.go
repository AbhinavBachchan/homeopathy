package handlers

import (
	"encoding/json"

	"homeopathy-platform/internal/models"
	"homeopathy-platform/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderHandler struct {
	db *gorm.DB
}

func NewOrderHandler(db *gorm.DB) *OrderHandler {
	return &OrderHandler{db: db}
}

type createOrderRequest struct {
	Items           []models.OrderItem `json:"items" binding:"required,min=1"`
	AddressID       uuid.UUID          `json:"address_id" binding:"required"`
	PaymentMethod   string              `json:"payment_method" binding:"required"`
	PrescriptionURL string              `json:"prescription_url"`
}

// Create builds an order from the cart. GST calculation, Schedule H
// prescription enforcement, and Razorpay/Stripe charge creation are the next
// pieces to wire in here before this is checkout-ready.
func (h *OrderHandler) Create(c *fiber.Ctx) error{
	var req createOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, 400, err.Error())
		
	}

	userIDStr := c.Locals("user_id")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		return response.Error(c, 401, "invalid user")
		
	}

	var subtotal int64
	var scheduleHPresent bool
	for _, item := range req.Items {
		var product models.Product
		if err := h.db.First(&product, "id = ?", item.ProductID).Error; err != nil {
			return response.Error(c, 404, "product not found: "+item.ProductID.String())
			
		}
		if product.Schedule == models.ScheduleH {
			scheduleHPresent = true
		}
		subtotal += product.Price * int64(item.Qty)
	}

	if scheduleHPresent && req.PrescriptionURL == "" {
		return response.Error(c, 400, "prescription required for Schedule H medicines in cart")
		
	}

	// GST: 12% on physical medicine products per brief section 10.1
	gst := subtotal * 12 / 100
	shipping := int64(0) // TODO: Shiprocket rate calc by pincode
	total := subtotal + gst + shipping

	itemsJSON, _ := json.Marshal(req.Items)

	order := models.Order{
		UserID:          userID,
		Status:          models.OrderPending,
		Items:           itemsJSON,
		Subtotal:        subtotal,
		GSTAmount:       gst,
		ShippingCharge:  shipping,
		Total:           total,
		PaymentMethod:   req.PaymentMethod,
		AddressID:       req.AddressID,
		PrescriptionURL: req.PrescriptionURL,
	}

	if err := h.db.Create(&order).Error; err != nil {
		return response.Error(c, 500, "failed to create order")
		
	}

	return response.Created(c, order)
	// TODO: kick off payment intent (Razorpay/Stripe) and return client secret
	// / order id for the frontend to complete payment.
}

func (h *OrderHandler) Get(c *fiber.Ctx) error{
	var order models.Order
	if err := h.db.First(&order, "id = ?", c.Params("id")).Error; err != nil {
		return response.Error(c, 404, "order not found")
		
	}
	return response.OK(c, order)
}

func (h *OrderHandler) ListMine(c *fiber.Ctx) error{
	userIDStr := c.Locals("user_id")
	var orders []models.Order
	if err := h.db.Where("user_id = ?", userIDStr).Order("created_at desc").Find(&orders).Error; err != nil {
		return response.Error(c, 500, "failed to fetch orders")
		
	}
	return response.OK(c, orders)
}
