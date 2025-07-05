package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"

	"ip-store/backend/internal/database"
	"ip-store/backend/internal/service"
)

// SetupRouter configures the routes for the application.
func SetupRouter(db *database.DBContext) *gin.Engine {
	r := gin.Default()

	// Configure CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Allow your Next.js frontend origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           3600, // 1 hour
	}))

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
