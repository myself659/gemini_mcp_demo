package api

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"ip-store/backend/internal/service"
)

func GetDownloadURLHandler(c *gin.Context) {
	orderID, err := strconv.ParseInt(c.Param("order_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	userID, _ := c.Get("userID")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	order, err := service.GetOrderByID(ctx, orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Verify if the user owns this order
	if order.UserID != userID.(int64) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to access this order"})
		return
	}

	// Generate the download URL for the product in the order
	downloadURL, err := service.GenerateDownloadURL(ctx, order.ProductID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"download_url": downloadURL})
}
