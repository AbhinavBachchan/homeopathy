package response

import "github.com/gofiber/fiber/v2"

func OK(c *fiber.Ctx, data any) error {
	return c.JSON(fiber.Map{"success": true, "data": data})
}

func Created(c *fiber.Ctx, data any) error {
	return c.Status(201).JSON(fiber.Map{"success": true, "data": data})
}

func Error(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{"success": false, "error": message})
}