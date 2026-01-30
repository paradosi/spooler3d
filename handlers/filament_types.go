package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/paradosi/spooler3d/models"
)

type FilamentTypeHandler struct {
	DB *sqlx.DB
}

func NewFilamentTypeHandler(db *sqlx.DB) *FilamentTypeHandler {
	return &FilamentTypeHandler{DB: db}
}

func (h *FilamentTypeHandler) List(c *gin.Context) {
	var items []models.FilamentType
	if err := h.DB.Select(&items, "SELECT * FROM filament_types ORDER BY name"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *FilamentTypeHandler) GetByID(c *gin.Context) {
	var item models.FilamentType
	if err := h.DB.Get(&item, "SELECT * FROM filament_types WHERE id = $1", c.Param("id")); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "filament type not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *FilamentTypeHandler) Create(c *gin.Context) {
	var req models.CreateFilamentTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var item models.FilamentType
	err := h.DB.Get(&item,
		`INSERT INTO filament_types (name, print_temp_min, print_temp_max, bed_temp_min, bed_temp_max)
		 VALUES ($1, $2, $3, $4, $5) RETURNING *`,
		req.Name, req.PrintTempMin, req.PrintTempMax, req.BedTempMin, req.BedTempMax)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *FilamentTypeHandler) Update(c *gin.Context) {
	var req models.CreateFilamentTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var item models.FilamentType
	err := h.DB.Get(&item,
		`UPDATE filament_types SET name = $1, print_temp_min = $2, print_temp_max = $3,
		 bed_temp_min = $4, bed_temp_max = $5, updated_at = $6 WHERE id = $7 RETURNING *`,
		req.Name, req.PrintTempMin, req.PrintTempMax, req.BedTempMin, req.BedTempMax, time.Now(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "filament type not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *FilamentTypeHandler) Delete(c *gin.Context) {
	result, err := h.DB.Exec("DELETE FROM filament_types WHERE id = $1", c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "filament type not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted": true})
}
