package handlers

import (
	"homeopathy-platform/internal/models"
	"homeopathy-platform/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
func (h *ProductHandler) List(c *gin.Context) {
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
		response.Error(c, 500, "failed to fetch products")
		return
	}

	response.OK(c, products)
}

func (h *ProductHandler) GetBySlug(c *gin.Context) {
	var product models.Product
	if err := h.db.Where("slug = ? AND is_active = ?", c.Param("slug"), true).First(&product).Error; err != nil {
		response.Error(c, 404, "product not found")
		return
	}
	response.OK(c, product)
}

// --- Admin CRUD (mount behind RequireRole(models.RoleAdmin)) ---

func (h *ProductHandler) Create(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	if err := h.db.Create(&product).Error; err != nil {
		response.Error(c, 500, "failed to create product")
		return
	}
	response.Created(c, product)
}

func (h *ProductHandler) Update(c *gin.Context) {
	var product models.Product
	if err := h.db.Where("id = ?", c.Param("id")).First(&product).Error; err != nil {
		response.Error(c, 404, "product not found")
		return
	}
	var updates models.Product
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	if err := h.db.Model(&product).Updates(updates).Error; err != nil {
		response.Error(c, 500, "failed to update product")
		return
	}
	response.OK(c, product)
}

func (h *ProductHandler) Delete(c *gin.Context) {
	if err := h.db.Where("id = ?", c.Param("id")).Delete(&models.Product{}).Error; err != nil {
		response.Error(c, 500, "failed to delete product")
		return
	}
	response.OK(c, gin.H{"deleted": true})
}

// TODO: bulk CSV import (P0 in brief) - stream-parse a CSV upload and
// batch-insert via h.db.CreateInBatches().
