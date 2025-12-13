package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type TransactionType string
type TransactionStatus string

const (
	TxDeposit       TransactionType = "DEPOSIT"
	TxWithdrawal    TransactionType = "WITHDRAWAL"
	TxTrade         TransactionType = "TRADE"
	TxFee           TransactionType = "FEE"
	TxTransfer      TransactionType = "TRANSFER"
	TxPromoRedeem   TransactionType = "PROMO_REDEEM"
	TxReferralBonus TransactionType = "REFERRAL_BONUS"

	TxStatusPending  TransactionStatus = "PENDING"
	TxStatusPosted   TransactionStatus = "POSTED"
	TxStatusFailed   TransactionStatus = "FAILED"
	TxStatusReversed TransactionStatus = "REVERSED"
)

type Balance struct {
	UserID      uuid.UUID       `gorm:"primaryKey;type:uuid" json:"user_id"`
	AssetSymbol string          `gorm:"primaryKey;size:20" json:"asset_symbol"`
	Available   decimal.Decimal `gorm:"type:numeric(30,18);not null;default:0" json:"available"`
	Locked      decimal.Decimal `gorm:"type:numeric(30,18);not null;default:0" json:"locked"`
	UpdatedAt   time.Time       `gorm:"autoUpdateTime" json:"updated_at"`

	// Gemini said, "Associations skipped for performance in Ledger"
	// Asset Asset `gorm:"foreignKey:AssetSymbol;references:Symbol"`
}

// Immutable history log.
type LedgerTransaction struct {
	ID           uuid.UUID         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	TxGroupID    uuid.UUID         `gorm:"type:uuid;index;not null" json:"tx_group_id"`
	UserID       uuid.UUID         `gorm:"type:uuid;index;not null" json:"user_id"`
	AssetSymbol  string            `gorm:"size:20;not null" json:"asset_symbol"`
	Amount       decimal.Decimal   `gorm:"type:numeric(30,18);not null" json:"amount"`
	Type         TransactionType   `gorm:"type:transaction_type;not null" json:"type"`
	Status       TransactionStatus `gorm:"type:transaction_status;not null;default:POSTED" json:"status"`
	ReferenceID  *uuid.UUID        `gorm:"type:uuid" json:"reference_id,omitempty"` // OrderID, WithdrawalID, etc.
	Description  string            `json:"description,omitempty"`
	BalanceAfter decimal.Decimal   `gorm:"type:numeric(30,18)" json:"balance_after"`
	CreatedAt    time.Time         `gorm:"autoCreateTime;index" json:"created_at"`
}
