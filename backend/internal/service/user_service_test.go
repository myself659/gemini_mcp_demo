package service

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *sql.DB {
	// Use in-memory database for testing
	db, err := sql.Open("sqlite", "file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("Failed to open in-memory database: %v", err)
	}

	// Create tables for testing
	createUserTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createUserTable)
	if err != nil {
		db.Close()
		t.Fatalf("Could not create users table for test: %v", err)
	}

	return db
}

func TestCreateUser(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	userService := NewUserService(db)

	// Test case 1: Successful user creation
	user, err := userService.CreateUser("test@example.com", "password123")
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}
	if user.ID == 0 {
		t.Error("Expected user ID to be non-zero")
	}
	if user.Email != "test@example.com" {
			t.Errorf("Expected email %s, got %s", "test@example.com", user.Email)
	}

	// Test case 2: User already exists
	_, err = userService.CreateUser("test@example.com", "anotherpassword")
	if err != ErrUserExists {
		t.Errorf("Expected ErrUserExists, got %v", err)
	}
}

func TestGetUserByEmail(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	userService := NewUserService(db)

	// Create a user first
	_, err := userService.CreateUser("findme@example.com", "securepassword")
	if err != nil {
		t.Fatalf("Failed to create user for test: %v", err)
	}

	// Test case 1: User found
	user, err := userService.GetUserByEmail("findme@example.com")
	if err != nil {
		t.Fatalf("GetUserByEmail failed: %v", err)
	}
	if user.Email != "findme@example.com" {
		t.Errorf("Expected email %s, got %s", "findme@example.com", user.Email)
	}

	// Test case 2: User not found
	_, err = userService.GetUserByEmail("nonexistent@example.com")
	if err != ErrUserNotFound {
		t.Errorf("Expected ErrUserNotFound, got %v", err)
	}
}

func TestComparePassword(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	userService := NewUserService(db)

	// Create a user to get a hashed password
	user, err := userService.CreateUser("passwordtest@example.com", "originalpassword")
	if err != nil {
		t.Fatalf("Failed to create user for password test: %v", err)
	}

	// Test case 1: Correct password
	err = userService.ComparePassword(user.PasswordHash, "originalpassword")
	if err != nil {
		t.Errorf("ComparePassword failed for correct password: %v", err)
	}

	// Test case 2: Incorrect password
	err = userService.ComparePassword(user.PasswordHash, "wrongpassword")
	if err == nil {
		t.Error("ComparePassword should have failed for incorrect password")
	}
}

func TestLoginUser(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	userService := NewUserService(db)

	// Create a user first
	_, err := userService.CreateUser("login@example.com", "loginpassword")
	if err != nil {
		t.Fatalf("Failed to create user for login test: %v", err)
	}

	// Test case 1: Successful login
	token, err := userService.LoginUser("login@example.com", "loginpassword")
	if err != nil {
		t.Fatalf("LoginUser failed: %v", err)
	}
	if token == "" {
		t.Error("Expected token to be non-empty")
	}

	// Test case 2: Incorrect password
	_, err = userService.LoginUser("login@example.com", "wrongpassword")
	if err == nil {
		t.Error("Expected error for incorrect password")
	}

	// Test case 3: User not found
	_, err = userService.LoginUser("nonexistent@example.com", "anypassword")
	if err != ErrUserNotFound {
		t.Errorf("Expected ErrUserNotFound, got %v", err)
	}
}
