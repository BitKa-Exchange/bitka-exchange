package app

import (
	"os"

	"bitka/pkg/config"
	"bitka/pkg/database"
	"bitka/pkg/middleware"
	"bitka/pkg/token"
	"bitka/services/account/internal/delivery/http"
	"bitka/services/account/internal/domain"
	"bitka/services/account/internal/repository"
	"bitka/services/account/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func NewServer(cfg *config.Config) (*fiber.App,domain.AccountUsecase, error) {
	// 1. Connect to ACCOUNT Database (Not Auth DB)
	db, err := database.Connect(database.Config{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPass,
		DBName:   cfg.DBName, // Different DB!
		SSLMode:  "disable",
	})
	if err != nil {
		return nil, nil,err
	}
	db.AutoMigrate(&domain.Profile{})

	// 2. Setup JWT Validator (Remote)
	jwksURL := os.Getenv("AUTH_JWKS_URL")
	if jwksURL == "" {
		// Fallback for local dev
		jwksURL = "http://localhost:3000/.well-known/jwks.json"
	}
	tokenValidator := token.NewValidator(jwksURL)

	// Create Middleware
	authMW := middleware.Protected(tokenValidator)

	// 3. DI
	repo := repository.NewAccountRepo(db)
	uc := usecase.NewAccountUsecase(repo)
	handler := http.NewAccountHandler(uc)

	// 4. Fiber
	app := fiber.New(fiber.Config{AppName: "Bitka Account Service"})
	app.Use(logger.New())
	app.Use(recover.New())

	http.MapRoutes(app, handler, authMW)

	return app, uc,nil
}
