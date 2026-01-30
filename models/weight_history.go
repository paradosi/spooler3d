package models

import "time"

type WeightHistory struct {
	ID         int       `db:"id" json:"id"`
	SpoolID    int       `db:"spool_id" json:"spool_id"`
	Weight     float64   `db:"weight" json:"weight"`
	MeasuredAt time.Time `db:"measured_at" json:"measured_at"`
}
