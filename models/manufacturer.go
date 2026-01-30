package models

import "time"

type Manufacturer struct {
	ID        int       `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Website   *string   `db:"website" json:"website"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type CreateManufacturerRequest struct {
	Name    string  `json:"name" binding:"required"`
	Website *string `json:"website"`
}
