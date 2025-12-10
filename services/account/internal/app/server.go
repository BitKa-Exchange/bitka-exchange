package app

import (
	"bitka/pkg/middleware"
	"bitka/services/account/internal/delivery/http"
	"bitka/services/account/internal/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"bitka/pkg/token"
)

func NewServer(uc domain.AccountUsecase, validator *token.Validator) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "Bitka Account Service",
	})

	app.Use(logger.New())
	app.Use(recover.New())

	authMW := middleware.Protected(validator)
	handler := http.NewAccountHandler(uc)

	http.MapRoutes(app, handler, authMW)

	return app
}
