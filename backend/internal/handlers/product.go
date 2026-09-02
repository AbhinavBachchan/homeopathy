package handlers

import (
	"homeopathy-platform/internal/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"homeopathy-platform/pkg/response"
)

type ProductHandler struct {
	db *gorm.DB
}

func NewProductHandler(db *gorm.DB) *ProductHandler {
	return &ProductHandler{db: db}
}

// List supports the filters called out in the brief: potency, category
// (therapeutic_category), brand (manufacturer), and free-text search.
// Algolia can be layered in later for instant/typo-tolerant search; this is
// the DB-backed fallback and what powers /api/products by default.
func (h *ProductHandler) List(c *fiber.Ctx) error{
	var products []models.Product
	query := h.db.Where("is_active = ?", true)

	if potency := c.Query("potency"); potency != "" {
		query = query.Where("potency = ?", potency)
	}
	if manufacturer := c.Query("brand"); manufacturer != "" {
		query = query.Where("manufacturer = ?", manufacturer)
	}
	if category := c.Query("category"); category != "" {
		query = query.Where("therapeutic_category = ?", category)
	}
	if q := c.Query("q"); q != "" {
		query = query.Where("name ILIKE ?", "%"+q+"%")
	}

	if err := query.Find(&products).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error":"failed to fetch products"})
	}

	return response.OK(c,products)
}

func (h *ProductHandler) GetBySlug(c *fiber.Ctx) error{
	var product models.Product
	if err := h.db.Where("slug = ? AND is_active = ?", c.Params("slug"), true).First(&product).Error; err != nil {
		return response.Error(c,404,"product not found")
	}
	return response.OK(c,product)
}

// --- Admin CRUD (mount behind RequireRole(models.RoleAdmin)) ---

func (h *ProductHandler) Create(c *fiber.Ctx) error{
	var product models.Product
	if err := c.BodyParser(&product); err != nil {
		return response.Error(c,400,err.Error())
	}
	if err := h.db.Create(&product).Error; err != nil {
		return response.Error(c,500,"failed to create product")
		
	}
	return response.Created(c,product)
}

func (h *ProductHandler) Update(c *fiber.Ctx) error{
	var product models.Product
	if err := h.db.Where("id = ?", c.Params("id")).First(&product).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"Error": "product not found"})
	}
	var updates models.Product
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if err := h.db.Model(&product).Updates(updates).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to update product"})
	}
	return c.JSON(fiber.Map{"product": product})
}

func (h *ProductHandler) Delete(c *fiber.Ctx) error{
	if err := h.db.Where("id = ?", c.Params("id")).Delete(&models.Product{}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to delete product"})
	}
	return c.JSON(fiber.Map{"deleted": true})
}

// TODO: bulk CSV import (P0 in brief) - stream-parse a CSV upload and
// batch-insert via h.db.CreateInBatches().
