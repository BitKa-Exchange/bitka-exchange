package http

import "github.com/gofiber/fiber/v2"

// MapRoutes now requires the JWT Middleware
func MapRoutes(app *fiber.App, h *AccountHandler, authMiddleware fiber.Handler) {
	api := app.Group("/v1/users")

	// Apply middleware to this group
	userGroup := api.Group("/users", authMiddleware)

	userGroup.Get("/me", h.GetProfile)
	userGroup.Put("/me", h.UpdateProfile)
}
