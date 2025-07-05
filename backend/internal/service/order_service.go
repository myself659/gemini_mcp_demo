package service

import (
	"ip-store/backend/internal/database"
	"ip-store/backend/internal/model"
)

func CreateOrder(order *model.Order) (int64, error) {
	return database.CreateOrder(order)
}

func GetOrderByID(id int64) (*model.Order, error) {
	return database.GetOrderByID(id)
}

func GetOrdersByUserID(userID int64) ([]*model.Order, error) {
	return database.GetOrdersByUserID(userID)
}

func ProcessPaymentWebhook(orderID int64, status string) error {
	// In a real application, you would add more robust validation here,
	// e.g., checking the payment gateway signature, verifying amount, etc.

	// For MVP, we simply update the order status.
	return database.UpdateOrderStatus(orderID, status)
}
