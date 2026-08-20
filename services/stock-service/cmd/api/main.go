package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/hottgbr/Korp_Teste_GabriellHott/services/stock-service/internal/config"
	"github.com/hottgbr/Korp_Teste_GabriellHott/services/stock-service/internal/database"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()

	db, err := database.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("could not start stock service: %v", err)
	}
	defer db.Close()

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "stock-service",
		})
	})

	log.Printf("stock service running on port %s", cfg.Port)

	if err := router.Run(cfg.Address()); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
}