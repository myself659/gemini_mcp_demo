package service

import (
	"context"

	"ip-store/backend/internal/database"
	"ip-store/backend/internal/model"
)

func CreateProduct(ctx context.Context, product *model.Product) (int64, error) {
	return database.CreateProduct(ctx, product)
}

func GetProductByID(ctx context.Context, id int64) (*model.Product, error) {
	return database.GetProductByID(ctx, id)
}

func GetAllProducts(ctx context.Context) ([]*model.Product, error) {
	return database.GetAllProducts(ctx)
}
