package service

import (
	"ip-store/backend/internal/database"
	"ip-store/backend/internal/model"
)

func CreateProduct(product *model.Product) (int64, error) {
	return database.CreateProduct(product)
}

func GetProductByID(id int64) (*model.Product, error) {
	return database.GetProductByID(id)
}

func GetAllProducts() ([]*model.Product, error) {
	return database.GetAllProducts()
}
