package main

import (
	"log"
	"os"

	"bitka/pkg/config"
	"bitka/pkg/database"
	"bitka/pkg/logger"
	"bitka/pkg/token"

	"bitka/services/account/internal/domain"
	"bitka/services/account/internal/repository"
	"bitka/services/account/internal/usecase"
	"bitka/services/account/internal/app"

	"bitka/services/account/internal/delivery/event"
)

func main() {
	logger.Init()
	cfg := config.Load("ACCOUNT_DB_NAME")

	if os.Getenv("DB_NAME") == "" {
		cfg.DBName = "bitka_account"
	}
	// 1. Connect DB (ONLY HERE)
	db, err := database.Connect(database.Config{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPass,
		DBName:   cfg.DBName,
		SSLMode:  "disable",
	})
	if err != nil {
		log.Fatalf("DB connect failed: %v", err)
	}
	db.AutoMigrate(&domain.Profile{})

	// 2. DI: Repo + Usecase
	repo := repository.NewAccountRepo(db)
	uc := usecase.NewAccountUsecase(repo)


	// 3. JWT Validator
	
	jwksURL := os.Getenv("AUTH_JWKS_URL")
	if jwksURL == "" {
		jwksURL = "http://localhost:3000/.well-known/jwks.json"
	}
	validator := token.NewValidator(jwksURL)

	// 4. Kafka Server

	kafkaHandler := event.NewHandler(uc)
	kafkaServer := event.NewServer(kafkaHandler)
	go kafkaServer.Start() 

	// 5. HTTP Server

	httpServer := app.NewServer(uc, validator)

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "3001"
	}

	log.Printf("Starting Account Service on :%s", port)
	httpServer.Listen(":" + port)
}
