package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	"ip-store/backend/internal/service"
)

// SetupRouter configures the routes for the application.
func SetupRouter(db *sql.DB) *gin.Engine {
	r := gin.Default()

	// Health check endpoint
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Initialize services and handlers
	userService := service.NewUserService(db)
	authHandler := NewAuthHandler(userService)

	// Auth routes
	authRoutes := r.Group("/api/auth")
	{
		authRoutes.POST("/register", authHandler.RegisterHandler)
		authRoutes.POST("/login", authHandler.LoginHandler)
	}

	// Product routes
	productRoutes := r.Group("/api/products")
	{
		productRoutes.GET("", GetProductsHandler)
		productRoutes.GET("/:id", GetProductHandler)
	}

	// Order routes
	orderRoutes := r.Group("/api/orders")
	orderRoutes.Use(AuthMiddleware())
	{
		orderRoutes.POST("", CreateOrderHandler)
		orderRoutes.GET("", GetOrdersHandler)
		orderRoutes.GET("/:id", GetOrderHandler)
	}

	// Payment webhook route
	paymentRoutes := r.Group("/api/payment")
	{
		paymentRoutes.POST("/webhook", PaymentWebhookHandler)
	}

	// Download routes
	downloadRoutes := r.Group("/api/downloads")
	downloadRoutes.Use(AuthMiddleware())
	{
		downloadRoutes.GET("/order/:order_id", GetDownloadURLHandler)
	}

	return r
}
