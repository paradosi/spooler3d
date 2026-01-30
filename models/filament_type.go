package models

import "time"

type FilamentType struct {
	ID           int       `db:"id" json:"id"`
	Name         string    `db:"name" json:"name"`
	PrintTempMin *int      `db:"print_temp_min" json:"print_temp_min"`
	PrintTempMax *int      `db:"print_temp_max" json:"print_temp_max"`
	BedTempMin   *int      `db:"bed_temp_min" json:"bed_temp_min"`
	BedTempMax   *int      `db:"bed_temp_max" json:"bed_temp_max"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type CreateFilamentTypeRequest struct {
	Name         string `json:"name" binding:"required"`
	PrintTempMin *int   `json:"print_temp_min"`
	PrintTempMax *int   `json:"print_temp_max"`
	BedTempMin   *int   `json:"bed_temp_min"`
	BedTempMax   *int   `json:"bed_temp_max"`
}
