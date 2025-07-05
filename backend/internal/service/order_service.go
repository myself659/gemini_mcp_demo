package service

import (
	"context"

	"ip-store/backend/internal/database"
	"ip-store/backend/internal/model"
)

func CreateOrder(ctx context.Context, order *model.Order) (int64, error) {
	return database.CreateOrder(ctx, order)
}

func GetOrderByID(ctx context.Context, id int64) (*model.Order, error) {
	return database.GetOrderByID(ctx, id)
}

func GetOrdersByUserID(ctx context.Context, userID int64) ([]*model.Order, error) {
	return database.GetOrdersByUserID(ctx, userID)
}

func ProcessPaymentWebhook(ctx context.Context, orderID int64, status string) error {
	// In a real application, you would add more robust validation here,
	// e.g., checking the payment gateway signature, verifying amount, etc.

	// For MVP, we simply update the order status.
	return database.UpdateOrderStatus(ctx, orderID, status)
}
