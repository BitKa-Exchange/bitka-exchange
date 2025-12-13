package dto

import "bitka/services/ledger/internal/domain"

type CreateAssetRequest struct {
	Symbol   string           `json:"symbol" validate:"required,alphanum,min=2,max=10"`
	Name     string           `json:"name" validate:"required"`
	Type     domain.AssetType `json:"type" validate:"required,oneof=CRYPTO FIAT"`
	Decimals int              `json:"decimals" validate:"required,min=0,max=18"`
	IconURL  string           `json:"icon_url" validate:"omitempty,url"`
	IsActive bool             `json:"is_active"`
}

type UpdateAssetRequest struct {
	Name     *string `json:"name,omitempty"`
	Decimals *int    `json:"decimals,omitempty" validate:"omitempty,min=0,max=18"`
	IconURL  *string `json:"icon_url,omitempty" validate:"omitempty,url"`
	IsActive *bool   `json:"is_active,omitempty"`
}

type AssetResponse struct {
	Symbol   string           `json:"symbol"`
	Name     string           `json:"name"`
	Type     domain.AssetType `json:"type"`
	Decimals int              `json:"decimals"`
	IconURL  string           `json:"icon_url"`
	IsActive bool             `json:"is_active"`
}
