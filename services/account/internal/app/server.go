package app

import (
	"bitka/pkg/config"
	"bitka/pkg/database"
	"bitka/pkg/token"
	"bitka/services/account/internal/delivery/event"
	"bitka/services/account/internal/delivery/http"
	"bitka/services/account/internal/domain"
	"bitka/services/account/internal/repository"
	"bitka/services/account/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	FiberServer *fiber.App
	KafkaServer *event.KafkaServer
}

func NewServer(cfg *config.Config) (*Server, error) {

	// 1. Connect DB
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

	db.AutoMigrate(&domain.Profile{})

	// 2. Token validator

	jwksURL := config.GetEnv("AUTH_JWKS_URL" , "http://localhost:3000/.well-known/jwks.json")
	
	validator := token.NewValidator(jwksURL)

	// 3. Construct usecase
	repo := repository.NewAccountRepo(db)
	uc := usecase.NewAccountUsecase(repo)

	// 4. Initialize Fiber
	httpServer := http.NewFiberServer(uc, validator)
	// 5. Initialize Kafka consumer (runs in background)
	kafkaconsumer := event.NewKafkaServer(uc)

	return &Server{
		FiberServer: httpServer,
		KafkaServer: kafkaconsumer,
	}, nil
}

func (s *Server) Listen(addr string) error {
	go s.KafkaServer.Start()
	return s.FiberServer.Listen(":8080")
}
