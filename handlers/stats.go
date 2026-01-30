package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type StatsHandler struct {
	DB *sqlx.DB
}

func NewStatsHandler(db *sqlx.DB) *StatsHandler {
	return &StatsHandler{DB: db}
}

type DashboardStats struct {
	TotalSpools      int      `db:"total_spools" json:"total_spools"`
	TotalWeight      *float64 `db:"total_weight" json:"total_weight"`
	TotalRemaining   *float64 `db:"total_remaining" json:"total_remaining"`
	ManufacturerCount int     `db:"manufacturer_count" json:"manufacturer_count"`
	TypeCount        int      `db:"type_count" json:"type_count"`
}

func (h *StatsHandler) Dashboard(c *gin.Context) {
	var stats DashboardStats

	err := h.DB.Get(&stats, `
		SELECT
			COUNT(*) AS total_spools,
			SUM(current_weight) AS total_weight,
			SUM(current_weight - spool_weight) AS total_remaining,
			(SELECT COUNT(*) FROM manufacturers) AS manufacturer_count,
			(SELECT COUNT(*) FROM filament_types) AS type_count
		FROM spools`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
