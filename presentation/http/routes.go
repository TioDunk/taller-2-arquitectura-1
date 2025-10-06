package http

import (
	"ddd-poc/presentation/http/handlers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configura todas las rutas de la aplicación
func SetupRoutes(
	customerHandler *handlers.CustomerHandler,
	productHandler *handlers.ProductHandler,
	orderHandler *handlers.OrderHandler,
) *gin.Engine {
	router := gin.Default()

	// Middleware para CORS
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Rutas de clientes
	customerRoutes := router.Group("/customers")
	{
		customerRoutes.POST("", customerHandler.CreateCustomer)
		customerRoutes.GET("/:id", customerHandler.GetCustomer)
		customerRoutes.GET("", customerHandler.GetAllCustomers)
		customerRoutes.PUT("/:id/address", customerHandler.UpdateCustomerAddress)
		customerRoutes.DELETE("/:id", customerHandler.DeleteCustomer)
	}

	// Rutas de productos
	productRoutes := router.Group("/products")
	{
		productRoutes.POST("", productHandler.CreateProduct)
		productRoutes.GET("/:id", productHandler.GetProduct)
		productRoutes.GET("", productHandler.GetAllProducts)
		productRoutes.GET("/active", productHandler.GetActiveProducts)
		productRoutes.PUT("/:id/price", productHandler.UpdateProductPrice)
		productRoutes.PUT("/:id/stock", productHandler.UpdateProductStock)
		productRoutes.PUT("/:id/deactivate", productHandler.DeactivateProduct)
	}

	// Rutas de pedidos
	orderRoutes := router.Group("/orders")
	{
		orderRoutes.POST("", orderHandler.CreateOrder)
		orderRoutes.GET("/:id", orderHandler.GetOrder)
		orderRoutes.GET("", orderHandler.GetAllOrders)
		orderRoutes.POST("/:id/add-item", orderHandler.AddItemToOrder)
		orderRoutes.PUT("/:id/confirm", orderHandler.ConfirmOrder)
		orderRoutes.PUT("/:id/status", orderHandler.UpdateOrderStatus)
	}

	// Ruta de salud
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "DDD POC API is running",
		})
	})

	return router
}
