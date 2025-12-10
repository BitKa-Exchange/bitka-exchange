package app

import (
	"log"
	"os"
	"bitka/pkg/config"
	"bitka/pkg/database"
	"bitka/pkg/token"
	"bitka/services/auth/internal/delivery/http"
	"bitka/services/auth/internal/domain"
	"bitka/services/auth/internal/delivery/event"
	"bitka/services/auth/internal/repository"
	"bitka/services/auth/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func NewServer(cfg *config.Config) (*fiber.App, error) {
	// 1. Infrastructure
	db, err := database.Connect(database.Config{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPass,
		DBName:   cfg.DBName,
		SSLMode:  "disable",
	})
	if err != nil {
		return nil, err
	}

	// Auto-Migrate auth tables
	db.AutoMigrate(&domain.User{}, &domain.RefreshToken{})

	// 2. Shared Components (Now using DB persistence)
	// We pass 'db' here so the manager can store keys in the database
	tokenMgr, err := token.NewManager(db)
	if err != nil {
		return nil, err
	}

	// 3. Layer Dependency Injection
	repo := repository.NewDatabaseRepo(db)
	broker := os.Getenv("KAFKA_BROKER")
	kafkaProducer, Err := event.NewProducer([]string{broker})
	if Err != nil {
		log.Fatal("Kafka producer failed:", Err)
	}
	uc := usecase.NewAuthUsecase(repo, tokenMgr, kafkaProducer)
	handler := http.NewAuthHandler(uc)

	// 4. Framework Setup
	app := fiber.New(fiber.Config{
		AppName: "Bitka Auth Service",
	})

	app.Use(recover.New())
	app.Use(logger.New())

	// 5. Route Mapping
	http.MapRoutes(app, handler)

	return app, nil
}
