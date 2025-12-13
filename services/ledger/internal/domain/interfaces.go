package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// LedgerRepository defines database operations (Postgres).
type LedgerRepository interface {
	// Asset Management
	CreateAsset(ctx context.Context, asset *Asset) error
	GetAsset(ctx context.Context, symbol string) (*Asset, error)
	ListAssets(ctx context.Context, activeOnly bool) ([]Asset, error)
	UpdateAsset(ctx context.Context, asset *Asset) error

	// Ledger Operations
	// TODO: updates Balances and inserts Transaction logs in one DB Transaction.
	CreateTransaction(ctx context.Context, tx LedgerTransaction) error

	GetBalance(ctx context.Context, userID uuid.UUID, assetSymbol string) (*Balance, error)
	ListBalances(ctx context.Context, userID uuid.UUID) ([]Balance, error)

	ListTransactions(ctx context.Context, filter TransactionFilter) ([]LedgerTransaction, int64, error)
}

// EventProducer defines message queue operations (Kafka).
// Treated as a Repository.
type EventProducer interface {
	PublishDepositConfirmed(ctx context.Context, tx LedgerTransaction) error
	PublishWithdrawalRequested(ctx context.Context, tx LedgerTransaction) error
}

// LedgerUsecase defines the business logic.
type LedgerUsecase interface {
	// Assets
	CreateAsset(ctx context.Context, req *Asset) error
	GetAssets(ctx context.Context, activeOnly bool) ([]Asset, error)
	GetAsset(ctx context.Context, symbol string) (*Asset, error)
	UpdateAsset(ctx context.Context, symbol string, req *Asset) error

	// Accounts
	GetUserBalances(ctx context.Context, userID uuid.UUID) ([]Balance, error)
	GetTransactions(ctx context.Context, filter TransactionFilter) ([]LedgerTransaction, int64, error)

	// Transactions (Logic)
	Deposit(ctx context.Context, userID uuid.UUID, asset string, amount decimal.Decimal) error
	Withdraw(ctx context.Context, userID uuid.UUID, asset string, amount decimal.Decimal) error
	Transfer(ctx context.Context, from, to uuid.UUID, asset string, amount decimal.Decimal) error
}

// TransactionFilter is a helper struct for listing queries
type TransactionFilter struct {
	UserID    uuid.UUID
	AccountID uuid.UUID // Optional alias for UserID in some contexts
	Asset     string
	Status    TransactionStatus
	From      time.Time
	To        time.Time
	Page      int
	PerPage   int
}
