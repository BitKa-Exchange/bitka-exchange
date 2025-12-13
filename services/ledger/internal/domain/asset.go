package domain

import (
	"time"
)

type AssetType string

const (
	AssetCrypto AssetType = "CRYPTO"
	AssetFiat   AssetType = "FIAT"
)

type Asset struct {
	Symbol    string    `gorm:"primaryKey;size:20" json:"symbol"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	Type      AssetType `gorm:"type:asset_type;not null" json:"type"`
	Decimals  int       `gorm:"not null;default:8" json:"decimals"`
	IconURL   string    `json:"icon_url"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
