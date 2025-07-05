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
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
	}

	return r
}
