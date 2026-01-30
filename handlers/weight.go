package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/paradosi/spooler3d/models"
)

type WeightHandler struct {
	DB *sqlx.DB
}

func NewWeightHandler(db *sqlx.DB) *WeightHandler {
	return &WeightHandler{DB: db}
}

type WeightUpdateRequest struct {
	Weight float64 `json:"weight" binding:"required"`
}

// UpdateWeight - POST /api/spools/:uid/weight
// The ESP32 sends the spool UUID (read from the NFC tag) and the measured weight.
func (h *WeightHandler) UpdateWeight(c *gin.Context) {
	uid := c.Param("uid")

	var req WeightUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := h.DB.Beginx()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to begin transaction"})
		return
	}
	defer tx.Rollback()

	now := time.Now()

	// Update the spool's current weight
	result, err := tx.Exec(
		"UPDATE spools SET current_weight = $1, updated_at = $2 WHERE uid = $3",
		req.Weight, now, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "spool not found"})
		return
	}

	// Record in weight history
	_, err = tx.Exec(
		`INSERT INTO weight_history (spool_id, weight, measured_at)
		 SELECT id, $1, $2 FROM spools WHERE uid = $3`,
		req.Weight, now, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated", "weight": req.Weight})
}

// GetHistory - GET /api/spools/:id/weight-history
func (h *WeightHandler) GetHistory(c *gin.Context) {
	var history []models.WeightHistory
	err := h.DB.Select(&history,
		"SELECT * FROM weight_history WHERE spool_id = $1 ORDER BY measured_at DESC",
		c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, history)
}
