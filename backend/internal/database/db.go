package database

import (
	"context"
	"database/sql"
	"log"
	"time"

	"ip-store/backend/internal/model"

	_ "modernc.org/sqlite"
)

// DBContext wraps the *sql.DB to provide context-aware database operations.
type DBContext struct {
	db *sql.DB
}

// NewDBContext creates a new DBContext instance.
func NewDBContext(db *sql.DB) *DBContext {
	return &DBContext{db: db}
}

// ExecContext executes a query without returning any rows.
func (ctx *DBContext) ExecContext(c context.Context, query string, args ...interface{}) (sql.Result, error) {
	return ctx.db.ExecContext(c, query, args...)
}

// QueryRowContext executes a query that is expected to return at most one row.
func (ctx *DBContext) QueryRowContext(c context.Context, query string, args ...interface{}) *sql.Row {
	return ctx.db.QueryRowContext(c, query, args...)
}

// QueryContext executes a query that returns rows, typically a SELECT statement.
func (ctx *DBContext) QueryContext(c context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return ctx.db.QueryContext(c, query, args...)
}

var DB *sql.DB
var DBConn *DBContext

func InitDB(dataSourceName string) {
	var err error
	// The driver name for modernc.org/sqlite is "sqlite"
	DB, err = sql.Open("sqlite", dataSourceName)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Set connection pool settings
	DB.SetMaxOpenConns(10) // Maximum number of open connections to the database
	DB.SetMaxIdleConns(5)  // Maximum number of connections in the idle connection pool
	DB.SetConnMaxLifetime(5 * time.Minute) // Maximum amount of time a connection may be reused

	if err = DB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	DBConn = NewDBContext(DB)

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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := DBConn.db.ExecContext(ctx, createUserTable)
	if err != nil {
		log.Fatalf("Could not create users table: %v", err)
	}

	_, err = DBConn.db.ExecContext(ctx, createProductTable)
	if err != nil {
		log.Fatalf("Could not create products table: %v", err)
	}

	_, err = DBConn.db.ExecContext(ctx, createOrderTable)
	if err != nil {
		log.Fatalf("Could not create orders table: %v", err)
	}
}

func GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}
	err := DBConn.db.QueryRowContext(ctx, "SELECT id, email, password_hash, created_at FROM users WHERE email = ?", email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func CreateProduct(ctx context.Context, product *model.Product) (int64, error) {
	result, err := DBConn.db.ExecContext(ctx, "INSERT INTO products (name, description, price, cover_image_url, file_key) VALUES (?, ?, ?, ?, ?)", product.Name, product.Description, product.Price, product.CoverImageURL, product.FileKey)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func GetProductByID(ctx context.Context, id int64) (*model.Product, error) {
	product := &model.Product{}
	err := DBConn.db.QueryRowContext(ctx, "SELECT id, name, description, price, cover_image_url, file_key, created_at FROM products WHERE id = ?", id).Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.CoverImageURL, &product.FileKey, &product.CreatedAt)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func GetAllProducts(ctx context.Context) ([]*model.Product, error) {
	rows, err := DBConn.db.QueryContext(ctx, "SELECT id, name, description, price, cover_image_url, file_key, created_at FROM products")
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

func CreateOrder(ctx context.Context, order *model.Order) (int64, error) {
	result, err := DBConn.db.ExecContext(ctx, "INSERT INTO orders (user_id, product_id, amount, status) VALUES (?, ?, ?, ?)", order.UserID, order.ProductID, order.Amount, order.Status)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func GetOrderByID(ctx context.Context, id int64) (*model.Order, error) {
	order := &model.Order{}
	err := DBConn.db.QueryRowContext(ctx, "SELECT id, user_id, product_id, amount, status, created_at, paid_at FROM orders WHERE id = ?", id).Scan(&order.ID, &order.UserID, &order.ProductID, &order.Amount, &order.Status, &order.CreatedAt, &order.PaidAt)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func GetOrdersByUserID(ctx context.Context, userID int64) ([]*model.Order, error) {
	rows, err := DBConn.db.QueryContext(ctx, "SELECT id, user_id, product_id, amount, status, created_at, paid_at FROM orders WHERE user_id = ?", userID)
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

func UpdateOrderStatus(ctx context.Context, orderID int64, status string) error {
	_, err := DBConn.db.ExecContext(ctx, "UPDATE orders SET status = ?, paid_at = CURRENT_TIMESTAMP WHERE id = ?", status, orderID)
	return err
}

func InsertInitialProducts(ctx context.Context, products []*model.Product) {
	for _, product := range products {
		// Check if product already exists by name
		var count int
		err := DBConn.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM products WHERE name = ?", product.Name).Scan(&count)
		if err != nil {
			log.Printf("Error checking product existence: %v", err)
			continue
		}
		if count > 0 {
			log.Printf("Product '%s' already exists, skipping insertion.", product.Name)
			continue
		}

		_, err = DBConn.db.ExecContext(ctx, "INSERT INTO products (name, description, price, cover_image_url, file_key) VALUES (?, ?, ?, ?, ?)",
			product.Name, product.Description, product.Price, product.CoverImageURL, product.FileKey)
		if err != nil {
			log.Printf("Error inserting product %s: %v", product.Name, err)
		}
	}
}
