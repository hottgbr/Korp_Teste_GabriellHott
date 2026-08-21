package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/hottgbr/Korp_Teste_GabriellHott/services/billing-service/internal/config"
	"github.com/hottgbr/Korp_Teste_GabriellHott/services/billing-service/internal/database"
	"github.com/hottgbr/Korp_Teste_GabriellHott/services/billing-service/internal/invoice"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()

	db, err := database.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf(
			"could not start billing service: %v",
			err,
		)
	}
	defer db.Close()
	invoiceRepository := invoice.NewRepository(db)
	invoiceService := invoice.NewService(invoiceRepository)
	invoiceHandler := invoice.NewHandler(invoiceService)
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "billing-service",
		})
	})
	invoices := router.Group("/invoices")
	{
		invoices.POST("", invoiceHandler.Create)
	}
	log.Printf(
		"billing service running on port %s",
		cfg.Port,
	)

	if err := router.Run(cfg.Address()); err != nil {
		log.Fatalf(
			"failed to start HTTP server: %v",
			err,
		)
	}
}
