package main

import (
	"fmt"
	"log"

	"ddd-poc/application/services"
	"ddd-poc/infrastructure/database"
	infraRepos "ddd-poc/infrastructure/repositories"
	presentationHttp "ddd-poc/presentation/http"
	"ddd-poc/presentation/http/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("🚀 Starting DDD POC Application...")

	// Crear la base de datos en memoria
	db := database.NewInMemoryDatabase()

	// Crear los repositorios
	customerRepo := infraRepos.NewInMemoryCustomerRepository(db)
	productRepo := infraRepos.NewInMemoryProductRepository(db)
	orderRepo := infraRepos.NewInMemoryOrderRepository(db)

	// Crear los servicios de aplicación
	customerService := services.NewCustomerService(customerRepo)
	productService := services.NewProductService(productRepo)
	orderService := services.NewOrderService(orderRepo, customerRepo, productRepo)

	// Crear los manejadores HTTP
	customerHandler := handlers.NewCustomerHandler(customerService)
	productHandler := handlers.NewProductHandler(productService)
	orderHandler := handlers.NewOrderHandler(orderService)

	// Configurar las rutas
	router := presentationHttp.SetupRoutes(customerHandler, productHandler, orderHandler)

	// Configurar Gin en modo release para producción
	gin.SetMode(gin.ReleaseMode)

	// Mostrar información de la aplicación
	fmt.Println("📋 Available endpoints:")
	fmt.Println("  GET    /health")
	fmt.Println("  POST   /customers")
	fmt.Println("  GET    /customers/:id")
	fmt.Println("  GET    /customers")
	fmt.Println("  PUT    /customers/:id/address")
	fmt.Println("  POST   /products")
	fmt.Println("  GET    /products/:id")
	fmt.Println("  GET    /products")
	fmt.Println("  GET    /products/active")
	fmt.Println("  PUT    /products/:id/price")
	fmt.Println("  PUT    /products/:id/stock")
	fmt.Println("  POST   /orders")
	fmt.Println("  GET    /orders/:id")
	fmt.Println("  GET    /orders")
	fmt.Println("  POST   /orders/:id/add-item")
	fmt.Println("  PUT    /orders/:id/confirm")
	fmt.Println("  PUT    /orders/:id/status")
	fmt.Println()
	fmt.Println("🌐 Server starting on port 8080...")
	fmt.Println("📖 Open your browser and go to: http://localhost:8080/health")
	fmt.Println()

	// Iniciar el servidor
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
