package http

import (
	"bitka/pkg/response"
	"bitka/services/auth/internal/delivery/http/dto"
	"bitka/services/auth/internal/domain"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	uc domain.AuthUsecase
}

func NewAuthHandler(uc domain.AuthUsecase) *AuthHandler {
	return &AuthHandler{uc: uc}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	tokens, err := h.uc.Login(req.Email, req.Password)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, err.Error())
	}

	return response.Success(c, dto.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.LoginRequest // Reusing for simplicity
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.uc.Register(req.Email, req.Password); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Registration failed")
	}

	return response.Success(c, "User registered successfully")
}

func (h *AuthHandler) GetJWKS(c *fiber.Ctx) error {
	keys, err := h.uc.GetJWKS()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get keys")
	}
	c.Set("Content-Type", "application/json")
	return c.Send(keys)
}
