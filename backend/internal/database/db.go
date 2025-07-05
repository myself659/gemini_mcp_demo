package database

import (
	"database/sql"
	"log"

	"ip-store/backend/internal/model"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB(dataSourceName string) {
	var err error
	// The driver name for modernc.org/sqlite is "sqlite"
	DB, err = sql.Open("sqlite", dataSourceName)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	createTables()
}

func createTables() {
	createUserTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	createProductTable := `
	CREATE TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(255) NOT NULL,
		description TEXT,
		price DECIMAL(10, 2) NOT NULL,
		cover_image_url VARCHAR(255),
		file_key VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	createOrderTable := `
	CREATE TABLE IF NOT EXISTS orders (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		product_id INTEGER NOT NULL,
		amount DECIMAL(10, 2) NOT NULL,
		status VARCHAR(50) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		paid_at TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (product_id) REFERENCES products(id)
	);`

	_, err := DB.Exec(createUserTable)
	if err != nil {
		log.Fatalf("Could not create users table: %v", err)
	}

	_, err = DB.Exec(createProductTable)
	if err != nil {
		log.Fatalf("Could not create products table: %v", err)
	}

	_, err = DB.Exec(createOrderTable)
	if err != nil {
		log.Fatalf("Could not create orders table: %v", err)
	}
}

func GetUserByEmail(email string) (*model.User, error) {
	user := &model.User{}
	err := DB.QueryRow("SELECT id, email, password_hash, created_at FROM users WHERE email = ?", email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func CreateProduct(product *model.Product) (int64, error) {
	result, err := DB.Exec("INSERT INTO products (name, description, price, cover_image_url, file_key) VALUES (?, ?, ?, ?, ?)", product.Name, product.Description, product.Price, product.CoverImageURL, product.FileKey)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func GetProductByID(id int64) (*model.Product, error) {
	product := &model.Product{}
	err := DB.QueryRow("SELECT id, name, description, price, cover_image_url, file_key, created_at FROM products WHERE id = ?", id).Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.CoverImageURL, &product.FileKey, &product.CreatedAt)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func GetAllProducts() ([]*model.Product, error) {
	rows, err := DB.Query("SELECT id, name, description, price, cover_image_url, file_key, created_at FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*model.Product
	for rows.Next() {
		product := &model.Product{}
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.CoverImageURL, &product.FileKey, &product.CreatedAt); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func CreateOrder(order *model.Order) (int64, error) {
	result, err := DB.Exec("INSERT INTO orders (user_id, product_id, amount, status) VALUES (?, ?, ?, ?)", order.UserID, order.ProductID, order.Amount, order.Status)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func GetOrderByID(id int64) (*model.Order, error) {
	order := &model.Order{}
	err := DB.QueryRow("SELECT id, user_id, product_id, amount, status, created_at, paid_at FROM orders WHERE id = ?", id).Scan(&order.ID, &order.UserID, &order.ProductID, &order.Amount, &order.Status, &order.CreatedAt, &order.PaidAt)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func GetOrdersByUserID(userID int64) ([]*model.Order, error) {
	rows, err := DB.Query("SELECT id, user_id, product_id, amount, status, created_at, paid_at FROM orders WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*model.Order
	for rows.Next() {
		order := &model.Order{}
		if err := rows.Scan(&order.ID, &order.UserID, &order.ProductID, &order.Amount, &order.Status, &order.CreatedAt, &order.PaidAt); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func UpdateOrderStatus(orderID int64, status string) error {
	_, err := DB.Exec("UPDATE orders SET status = ?, paid_at = CURRENT_TIMESTAMP WHERE id = ?", status, orderID)
	return err
}
