package main

import (
	"github.com/raffaelramalhorosa/empacotador-api/handlers"

	"github.com/gin-gonic/gin"
	_ "github.com/raffaelramalhorosa/empacotador-api/docs" // Documentação gerada pelo Swagger
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title API de Empacotamento de Pedidos
// @version 1.2
// @description API para otimizar o empacotamento de produtos em caixas de papelão
// @host localhost:8080
// @BasePath /
func main() {
	// Inicializa Gin
	r := gin.Default()

	// Handler
	orderHandler := handlers.NewOrderHandler()

	// Rotas
	r.POST("/api/pack", orderHandler.PackOrders)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Rota de health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Inicia servidor
	r.Run(":8080")
}
