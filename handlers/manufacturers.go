package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/paradosi/spooler3d/models"
)

type ManufacturerHandler struct {
	DB *sqlx.DB
}

func NewManufacturerHandler(db *sqlx.DB) *ManufacturerHandler {
	return &ManufacturerHandler{DB: db}
}

func (h *ManufacturerHandler) List(c *gin.Context) {
	var items []models.Manufacturer
	if err := h.DB.Select(&items, "SELECT * FROM manufacturers ORDER BY name"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *ManufacturerHandler) GetByID(c *gin.Context) {
	var item models.Manufacturer
	if err := h.DB.Get(&item, "SELECT * FROM manufacturers WHERE id = $1", c.Param("id")); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "manufacturer not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *ManufacturerHandler) Create(c *gin.Context) {
	var req models.CreateManufacturerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var item models.Manufacturer
	err := h.DB.Get(&item,
		"INSERT INTO manufacturers (name, website) VALUES ($1, $2) RETURNING *",
		req.Name, req.Website)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *ManufacturerHandler) Update(c *gin.Context) {
	var req models.CreateManufacturerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var item models.Manufacturer
	err := h.DB.Get(&item,
		"UPDATE manufacturers SET name = $1, website = $2, updated_at = $3 WHERE id = $4 RETURNING *",
		req.Name, req.Website, time.Now(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "manufacturer not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *ManufacturerHandler) Delete(c *gin.Context) {
	result, err := h.DB.Exec("DELETE FROM manufacturers WHERE id = $1", c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "manufacturer not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted": true})
}
