package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/paradosi/spooler3d/models"
)

func generateUID() string {
	b := make([]byte, 7)
	rand.Read(b)
	return hex.EncodeToString(b)
}

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

	uid := ""
	if req.UID != nil && *req.UID != "" {
		uid = *req.UID
	} else {
		uid = generateUID()
	}

	var spool models.Spool
	err := h.DB.Get(&spool,
		`INSERT INTO spools
		 (uid, manufacturer_id, filament_type_id, color_name, color_hex,
		  diameter, spool_weight, net_weight, current_weight,
		  location, purchase_date, purchase_price, notes, td_code)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		 RETURNING *, (current_weight - spool_weight) AS remaining_weight`,
		uid, req.ManufacturerID, req.FilamentTypeID, req.ColorName, req.ColorHex,
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
		 uid = COALESCE($1, uid),
		 manufacturer_id = $2, filament_type_id = $3, color_name = $4, color_hex = $5,
		 diameter = $6, spool_weight = $7, net_weight = $8, current_weight = $9,
		 location = $10, purchase_date = $11, purchase_price = $12, notes = $13,
		 td_code = $14, updated_at = $15
		 WHERE id = $16
		 RETURNING *, (current_weight - spool_weight) AS remaining_weight`,
		req.UID, req.ManufacturerID, req.FilamentTypeID, req.ColorName, req.ColorHex,
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
