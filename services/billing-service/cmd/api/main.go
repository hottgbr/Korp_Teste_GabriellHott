package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/hottgbr/Korp_Teste_GabriellHott/services/billing-service/internal/config"
	"github.com/hottgbr/Korp_Teste_GabriellHott/services/billing-service/internal/database"
	"github.com/hottgbr/Korp_Teste_GabriellHott/services/billing-service/internal/invoice"
	"github.com/hottgbr/Korp_Teste_GabriellHott/services/billing-service/internal/stockclient"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()

	db, err := database.NewPostgresPool(
		ctx,
		cfg.DatabaseURL,
	)
	if err != nil {
		log.Fatalf(
			"could not start billing service: %v",
			err,
		)
	}
	defer db.Close()

	stockClient := stockclient.NewClient(
		cfg.StockServiceURL,
	)

	invoiceRepository := invoice.NewRepository(db)

	invoiceService := invoice.NewService(
		invoiceRepository,
		stockClient,
	)

	invoiceHandler := invoice.NewHandler(
		invoiceService,
	)

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "billing-service",
		})
	})

	invoices := router.Group("/invoices")
	{
		invoices.POST(
			"",
			invoiceHandler.Create,
		)

		invoices.GET(
			"",
			invoiceHandler.List,
		)

		invoices.GET(
			"/:id",
			invoiceHandler.FindByID,
		)

		invoices.POST(
			"/:id/close",
			invoiceHandler.Close,
		)
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
