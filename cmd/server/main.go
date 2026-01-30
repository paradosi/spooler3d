package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/paradosi/spooler3d/config"
	"github.com/paradosi/spooler3d/db"
	"github.com/paradosi/spooler3d/router"
)

func main() {
	cfg := config.Load()
	gin.SetMode(cfg.GinMode)

	database, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer database.Close()

	log.Println("Connected to PostgreSQL")

	r := router.Setup(database)

	log.Printf("Spooler3D API starting on :%s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
