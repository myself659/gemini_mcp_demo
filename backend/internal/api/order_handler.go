package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"ip-store/backend/internal/model"
	"ip-store/backend/internal/service"
)

func CreateOrderHandler(c *gin.Context) {
	var order model.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")
	order.UserID = userID.(int64)

	orderID, err := service.CreateOrder(&order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"order_id": orderID})
}

func GetOrderHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	order, err := service.GetOrderByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	userID, _ := c.Get("userID")
	if order.UserID != userID.(int64) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to view this order"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func GetOrdersHandler(c *gin.Context) {
	userID, _ := c.Get("userID")

	orders, err := service.GetOrdersByUserID(userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}
