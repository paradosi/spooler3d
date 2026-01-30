package models

import (
	"time"
)

type Spool struct {
	ID              int        `db:"id" json:"id"`
	UID             string     `db:"uid" json:"uid"`
	ManufacturerID  *int       `db:"manufacturer_id" json:"manufacturer_id"`
	FilamentTypeID  *int       `db:"filament_type_id" json:"filament_type_id"`
	ColorName       *string    `db:"color_name" json:"color_name"`
	ColorHex        *string    `db:"color_hex" json:"color_hex"`
	Diameter        float64    `db:"diameter" json:"diameter"`
	SpoolWeight     *float64   `db:"spool_weight" json:"spool_weight"`
	NetWeight       *float64   `db:"net_weight" json:"net_weight"`
	CurrentWeight   *float64   `db:"current_weight" json:"current_weight"`
	RemainingWeight *float64   `db:"remaining_weight" json:"remaining_weight"`
	Location        *string    `db:"location" json:"location"`
	PurchaseDate    *time.Time `db:"purchase_date" json:"purchase_date"`
	PurchasePrice   *float64   `db:"purchase_price" json:"purchase_price"`
	Notes           *string    `db:"notes" json:"notes"`
	TDCode          *string    `db:"td_code" json:"td_code"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at" json:"updated_at"`
}

type CreateSpoolRequest struct {
	UID            *string  `json:"uid"`
	ManufacturerID *int     `json:"manufacturer_id"`
	FilamentTypeID *int     `json:"filament_type_id"`
	ColorName      *string  `json:"color_name"`
	ColorHex       *string  `json:"color_hex"`
	Diameter       float64  `json:"diameter" binding:"required"`
	SpoolWeight    *float64 `json:"spool_weight"`
	NetWeight      *float64 `json:"net_weight"`
	CurrentWeight  *float64 `json:"current_weight"`
	Location       *string  `json:"location"`
	PurchaseDate   *string  `json:"purchase_date"`
	PurchasePrice  *float64 `json:"purchase_price"`
	Notes          *string  `json:"notes"`
	TDCode         *string  `json:"td_code"`
}
