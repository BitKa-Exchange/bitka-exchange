package http

import "github.com/gofiber/fiber/v2"

func MapRoutes(app *fiber.App, h *AuthHandler) {
	api := app.Group("/api/v1")

	api.Post("/login", h.Login)
	api.Post("/register", h.Register)

	// JWKS endpoint often lives at root or .well-known
	app.Get("/.well-known/jwks.json", h.GetJWKS)
}
