package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"ip-store/backend/internal/service"
)

type PaymentWebhookRequest struct {
	OrderID int64  `json:"order_id" binding:"required"`
	Status  string `json:"status" binding:"required"`
}

func PaymentWebhookHandler(c *gin.Context) {
	var req PaymentWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// In a real application, you would verify the webhook signature here
	// to ensure it's from a legitimate payment gateway.

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	err := service.ProcessPaymentWebhook(ctx, req.OrderID, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook processed successfully"})
}
