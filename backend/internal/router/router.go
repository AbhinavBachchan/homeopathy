package router

import (
	"homeopathy-platform/internal/config"
	"homeopathy-platform/internal/handlers"
	"homeopathy-platform/internal/middleware"
	"homeopathy-platform/internal/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func New(cfg *config.Config, db *gorm.DB) *fiber.App {
	app := fiber.New()

	// CORS for the Angular dev server / production frontend origin.
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*") // tighten to your Angular origin in production
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		if c.Method() == "OPTIONS" {
			return c.SendStatus(204)
		}
		return c.Next()
	})

	authHandler := handlers.NewAuthHandler(db, cfg)
	productHandler := handlers.NewProductHandler(db)
	orderHandler := handlers.NewOrderHandler(db)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	api := app.Group("/api")
	{
		// Auth
		api.Post("/auth/register", authHandler.Register)
		api.Post("/auth/login", authHandler.Login)
		api.Post("/auth/forgot-password", authHandler.ForgotPassword)
		api.Post("/auth/reset-password", authHandler.ResetPassword)

		// Products (public read)
		api.Get("/products", productHandler.List)
		api.Get("/products/:slug", productHandler.GetBySlug)

		// Orders (patient, authenticated)
		authed := api.Group("")
		authed.Use(middleware.RequireAuth(cfg))
		authed.Post("/orders", orderHandler.Create)
		authed.Get("/orders/:id", orderHandler.Get)
		authed.Get("/orders", orderHandler.ListMine)

		// Admin (product CRUD, dashboard)
		admin := api.Group("/admin")
		admin.Use(middleware.RequireAuth(cfg), middleware.RequireRole(models.RoleAdmin))
		admin.Post("/products", productHandler.Create)
		admin.Put("/products/:id", productHandler.Update)
		admin.Delete("/products/:id", productHandler.Delete)

		// TODO next: /api/webhooks/razorpay, /api/webhooks/stripe,
		// /api/webhooks/shiprocket, /api/webhooks/interakt (unauthenticated,
		// signature-verified) and /api/ai/symptom-assess once Phase 2 starts.
	}
	return app
}
