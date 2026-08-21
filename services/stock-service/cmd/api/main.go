package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/hottgbr/Korp_Teste_GabriellHott/services/stock-service/internal/config"
	"github.com/hottgbr/Korp_Teste_GabriellHott/services/stock-service/internal/database"
	"github.com/hottgbr/Korp_Teste_GabriellHott/services/stock-service/internal/product"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()

	db, err := database.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("could not start stock service: %v", err)
	}
	defer db.Close()
	productRepository := product.NewRepository(db)
	productService := product.NewService(productRepository)
	productHandler := product.NewHandler(productService)
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "stock-service",
		})
	})
	products := router.Group("/products")
	{
		products.POST("", productHandler.Create)
		products.GET("", productHandler.List)
		products.GET("/:id", productHandler.FindByID)
	}
	log.Printf("stock service running on port %s", cfg.Port)

	if err := router.Run(cfg.Address()); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
}
