package http

import (
	"bitka/pkg/middleware"
	"bitka/pkg/token"
	"bitka/services/account/internal/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func NewFiberServer(uc domain.AccountUsecase, validator *token.Validator) *fiber.App {
	FiberServer := fiber.New()

	FiberServer.Use(logger.New())
	FiberServer.Use(recover.New())

	authMW := middleware.Protected(validator)
	handler := NewAccountHandler(uc)

	MapRoutes(FiberServer, handler, authMW)

	return FiberServer
}
