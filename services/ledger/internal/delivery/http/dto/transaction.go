package dto

import (
	"time"

	"bitka/services/ledger/internal/domain"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreateTransactionRequest struct {
	AccountID   uuid.UUID       `json:"account_id" validate:"required"` // UserID
	Type        string          `json:"type" validate:"required,oneof=debit credit transfer"`
	Amount      decimal.Decimal `json:"amount" validate:"required"`
	Asset       string          `json:"asset" validate:"required"`
	Description string          `json:"description"`
	// Metadata maps string->interface{} in Swagger, usually for flexible contexts
	Metadata map[string]interface{} `json:"metadata"`
}

type TransactionResponse struct {
	ID           uuid.UUID                `json:"id"`
	UserID       uuid.UUID                `json:"account_id"` // Mapped to account_id for API consistency
	Type         domain.TransactionType   `json:"type"`
	Amount       decimal.Decimal          `json:"amount"`
	Asset        string                   `json:"asset"`
	Status       domain.TransactionStatus `json:"status"`
	Description  string                   `json:"description,omitempty"`
	BalanceAfter decimal.Decimal          `json:"balance_after"`
	CreatedAt    time.Time                `json:"created_at"`
}

type CreateTransactionResponse struct {
	TransactionID uuid.UUID `json:"transaction_id"`
	Status        string    `json:"status"`
	Message       string    `json:"message"`
}

type BalanceResponse struct {
	UserID    uuid.UUID       `json:"user_id"`
	Asset     string          `json:"asset"`
	Available decimal.Decimal `json:"available"`
	Locked    decimal.Decimal `json:"locked"`
}
