package http

import (
	"bitka/pkg/response"
	"bitka/services/account/internal/delivery/http/dto"
	"bitka/services/account/internal/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AccountHandler struct {
	uc domain.AccountUsecase
}

func NewAccountHandler(uc domain.AccountUsecase) *AccountHandler {
	return &AccountHandler{uc: uc}
}

func (h *AccountHandler) GetProfile(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	profile, err := h.uc.GetMyProfile(userID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, profile)
}

func (h *AccountHandler) UpdateProfile(c *fiber.Ctx) error {
	var req dto.UpdateProfileRequest
	userIDStr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request")
	}

	err = h.uc.UpdateMyProfile(userID, req.FullName, req.AvatarURL)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, "Profile updated")
}
