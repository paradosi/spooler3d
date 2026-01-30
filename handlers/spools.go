package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/paradosi/spooler3d/models"
)

type SpoolHandler struct {
	DB *sqlx.DB
}

func NewSpoolHandler(db *sqlx.DB) *SpoolHandler {
	return &SpoolHandler{DB: db}
}

const spoolSelectQuery = `SELECT *,
	(current_weight - spool_weight) AS remaining_weight
	FROM spools`

func (h *SpoolHandler) List(c *gin.Context) {
	var spools []models.Spool

	query := spoolSelectQuery + " ORDER BY updated_at DESC"

	// Optional filters
	mfr := c.Query("manufacturer_id")
	ft := c.Query("filament_type_id")
	loc := c.Query("location")

	if mfr != "" {
		query = spoolSelectQuery + " WHERE manufacturer_id = $1 ORDER BY updated_at DESC"
		if err := h.DB.Select(&spools, query, mfr); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, spools)
		return
	}
	if ft != "" {
		query = spoolSelectQuery + " WHERE filament_type_id = $1 ORDER BY updated_at DESC"
		if err := h.DB.Select(&spools, query, ft); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, spools)
		return
	}
	if loc != "" {
		query = spoolSelectQuery + " WHERE location ILIKE $1 ORDER BY updated_at DESC"
		if err := h.DB.Select(&spools, query, "%"+loc+"%"); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, spools)
		return
	}

	if err := h.DB.Select(&spools, query); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, spools)
}

func (h *SpoolHandler) GetByID(c *gin.Context) {
	var spool models.Spool
	query := spoolSelectQuery + " WHERE id = $1"
	if err := h.DB.Get(&spool, query, c.Param("id")); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "spool not found"})
		return
	}
	c.JSON(http.StatusOK, spool)
}

func (h *SpoolHandler) Create(c *gin.Context) {
	var req models.CreateSpoolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var spool models.Spool
	err := h.DB.Get(&spool,
		`INSERT INTO spools
		 (manufacturer_id, filament_type_id, color_name, color_hex,
		  diameter, spool_weight, net_weight, current_weight,
		  location, purchase_date, purchase_price, notes, td_code)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		 RETURNING *, (current_weight - spool_weight) AS remaining_weight`,
		req.ManufacturerID, req.FilamentTypeID, req.ColorName, req.ColorHex,
		req.Diameter, req.SpoolWeight, req.NetWeight, req.CurrentWeight,
		req.Location, req.PurchaseDate, req.PurchasePrice, req.Notes, req.TDCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, spool)
}

func (h *SpoolHandler) Update(c *gin.Context) {
	var req models.CreateSpoolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var spool models.Spool
	err := h.DB.Get(&spool,
		`UPDATE spools SET
		 manufacturer_id = $1, filament_type_id = $2, color_name = $3, color_hex = $4,
		 diameter = $5, spool_weight = $6, net_weight = $7, current_weight = $8,
		 location = $9, purchase_date = $10, purchase_price = $11, notes = $12,
		 td_code = $13, updated_at = $14
		 WHERE id = $15
		 RETURNING *, (current_weight - spool_weight) AS remaining_weight`,
		req.ManufacturerID, req.FilamentTypeID, req.ColorName, req.ColorHex,
		req.Diameter, req.SpoolWeight, req.NetWeight, req.CurrentWeight,
		req.Location, req.PurchaseDate, req.PurchasePrice, req.Notes,
		req.TDCode, time.Now(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "spool not found"})
		return
	}
	c.JSON(http.StatusOK, spool)
}

func (h *SpoolHandler) Delete(c *gin.Context) {
	result, err := h.DB.Exec("DELETE FROM spools WHERE id = $1", c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "spool not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted": true})
}
