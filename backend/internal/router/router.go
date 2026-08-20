package router

import (
	"homeopathy-platform/internal/config"
	"homeopathy-platform/internal/handlers"
	"homeopathy-platform/internal/middleware"
	"homeopathy-platform/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func New(cfg *config.Config, db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// CORS for the Angular dev server / production frontend origin.
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*") // tighten to your Angular origin in production
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	authHandler := handlers.NewAuthHandler(db, cfg)
	productHandler := handlers.NewProductHandler(db)
	orderHandler := handlers.NewOrderHandler(db)

	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	api := r.Group("/api")
	{
		// Auth
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)

		// Products (public read)
		api.GET("/products", productHandler.List)
		api.GET("/products/:slug", productHandler.GetBySlug)

		// Orders (patient, authenticated)
		authed := api.Group("")
		authed.Use(middleware.RequireAuth(cfg))
		{
			authed.POST("/orders", orderHandler.Create)
			authed.GET("/orders/:id", orderHandler.Get)
			authed.GET("/orders", orderHandler.ListMine)
		}

		// Admin (product CRUD, dashboard)
		admin := api.Group("/admin")
		admin.Use(middleware.RequireAuth(cfg), middleware.RequireRole(models.RoleAdmin))
		{
			admin.POST("/products", productHandler.Create)
			admin.PUT("/products/:id", productHandler.Update)
			admin.DELETE("/products/:id", productHandler.Delete)
		}

		// TODO next: /api/webhooks/razorpay, /api/webhooks/stripe,
		// /api/webhooks/shiprocket, /api/webhooks/interakt (unauthenticated,
		// signature-verified) and /api/ai/symptom-assess once Phase 2 starts.
	}

	return r
}
