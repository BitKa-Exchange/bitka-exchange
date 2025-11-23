package http

import (
	"bitka/pkg/response"
	"bitka/services/account/internal/domain"

	"github.com/gofiber/fiber/v2"
)

type AccountHandler struct {
	uc domain.AccountUsecase
}

func NewAccountHandler(uc domain.AccountUsecase) *AccountHandler {
	return &AccountHandler{uc: uc}
}

func (h *AccountHandler) GetProfile(c *fiber.Ctx) error {
	// "user_id" is set by the Middleware
	userID := c.Locals("user_id").(string)

	profile, err := h.uc.GetMyProfile(userID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, profile)
}

func (h *AccountHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var req struct {
		FullName  string `json:"full_name"`
		AvatarURL string `json:"avatar_url"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request")
	}

	err := h.uc.UpdateMyProfile(userID, req.FullName, req.AvatarURL)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, "Profile updated")
}

// MapRoutes now requires the JWT Middleware
func MapRoutes(app *fiber.App, h *AccountHandler, authMiddleware fiber.Handler) {
	api := app.Group("/api/v1")

	// Apply middleware to this group
	userGroup := api.Group("/users", authMiddleware)

	userGroup.Get("/me", h.GetProfile)
	userGroup.Put("/me", h.UpdateProfile)
}
