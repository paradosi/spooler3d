package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/paradosi/spooler3d/handlers"
)

func Setup(db *sqlx.DB) *gin.Engine {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	api := r.Group("/api")

	// Health
	api.GET("/health", handlers.HealthCheck(db))

	// Manufacturers
	mh := handlers.NewManufacturerHandler(db)
	api.GET("/manufacturers", mh.List)
	api.GET("/manufacturers/:id", mh.GetByID)
	api.POST("/manufacturers", mh.Create)
	api.PUT("/manufacturers/:id", mh.Update)
	api.DELETE("/manufacturers/:id", mh.Delete)

	// Filament Types
	fth := handlers.NewFilamentTypeHandler(db)
	api.GET("/filament-types", fth.List)
	api.GET("/filament-types/:id", fth.GetByID)
	api.POST("/filament-types", fth.Create)
	api.PUT("/filament-types/:id", fth.Update)
	api.DELETE("/filament-types/:id", fth.Delete)

	// Spools
	sh := handlers.NewSpoolHandler(db)
	api.GET("/spools", sh.List)
	api.GET("/spools/:id", sh.GetByID)
	api.POST("/spools", sh.Create)
	api.PUT("/spools/:id", sh.Update)
	api.DELETE("/spools/:id", sh.Delete)

	// Weight (ESP32 uses UUID from NFC tag)
	wh := handlers.NewWeightHandler(db)
	api.POST("/spools/:uid/weight", wh.UpdateWeight)
	api.GET("/spools/:id/weight-history", wh.GetHistory)

	// Stats
	stath := handlers.NewStatsHandler(db)
	api.GET("/stats", stath.Dashboard)

	return r
}
