package handlers

import (
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

type createOrderItemRequest struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	Qty       int       `json:"qty" binding:"required,min=1"`
}

type createOrderRequest struct {
	Items           []createOrderItemRequest `json:"items" binding:"required,min=1"`
	AddressID       uuid.UUID                `json:"address_id" binding:"required"`
	PaymentMethod   string                   `json:"payment_method" binding:"required"`
	PrescriptionURL string                   `json:"prescription_url"`
}

// Create creates an order and its associated order items.
func (h *OrderHandler) Create(c *fiber.Ctx) error {
	var req createOrderRequest

	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, 400, err.Error())
	}

	userIDStr := c.Locals("user_id")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		return response.Error(c, 401, "invalid user")
	}

	// Verify that the address exists and belongs to the current user.
	var address models.Address
	if err := h.db.Where("id = ? AND user_id = ?", req.AddressID, userID).First(&address).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.Error(c, 404, "address not found")
		}

		return response.Error(c, 500, "failed to validate address")
	}

	var subtotal int64
	var scheduleHPresent bool

	// These are the items we will actually save.
	// Name, potency, and price come from the database product.
	orderItems := make([]models.OrderItem, 0, len(req.Items))

	for _, reqItem := range req.Items {
		if reqItem.Qty <= 0 {
			return response.Error(c, 400, "quantity must be greater than 0")
		}

		var product models.Product

		if err := h.db.
			Where("id = ? AND is_active = ?", reqItem.ProductID, true).
			First(&product).Error; err != nil {

			if err == gorm.ErrRecordNotFound {
				return response.Error(
					c,
					404,
					"product not found: "+reqItem.ProductID.String(),
				)
			}

			return response.Error(c, 500, "failed to fetch product")
		}

		if product.StockQty < reqItem.Qty {
			return response.Error(
				c,
				400,
				"insufficient stock for product: "+product.Name,
			)
		}

		if product.Schedule == models.ScheduleH {
			scheduleHPresent = true
		}

		subtotal += product.Price * int64(reqItem.Qty)

		orderItems = append(orderItems, models.OrderItem{
			ProductID: product.ID,
			Name:      product.Name,
			Potency:   product.Potency,
			Qty:       reqItem.Qty,
			Price:     product.Price,
		})
	}

	if scheduleHPresent && req.PrescriptionURL == "" {
		return response.Error(
			c,
			400,
			"prescription required for Schedule H medicines in cart",
		)
	}

	// GST: 12%
	gst := subtotal * 12 / 100
	shipping := int64(0)
	total := subtotal + gst + shipping

	var order models.Order

	// Create the order and order items atomically.
	err = h.db.Transaction(func(tx *gorm.DB) error {
		order = models.Order{
			UserID:          userID,
			Status:          models.OrderPending,
			Subtotal:        subtotal,
			GSTAmount:       gst,
			ShippingCharge:  shipping,
			Total:           total,
			PaymentMethod:   req.PaymentMethod,
			AddressID:       req.AddressID,
			PrescriptionURL: req.PrescriptionURL,
		}

		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		for i := range orderItems {
			orderItems[i].OrderID = order.ID
		}

		if err := tx.Create(&orderItems).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return response.Error(c, 500, "failed to create order")
	}

	// Attach items for the response.
	order.Items = orderItems

	return response.Created(c, order)

}

func (h *OrderHandler) Get(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id")

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		return response.Error(c, 401, "invalid user")
	}

	var order models.Order

	if err := h.db.
		Preload("Items").
		Where("id = ? AND user_id = ?", c.Params("id"), userID).
		First(&order).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return response.Error(c, 404, "order not found")
		}

		return response.Error(c, 500, "failed to fetch order")
	}

	return response.OK(c, order)

}

func (h *OrderHandler) ListMine(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id")

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		return response.Error(c, 401, "invalid user")
	}

	var orders []models.Order

	if err := h.db.
		Where("user_id = ?", userID).
		Order("created_at desc").
		Preload("Items").
		Find(&orders).Error; err != nil {

		return response.Error(c, 500, "failed to fetch orders")
	}

	return response.OK(c, orders)

}
