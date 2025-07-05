package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"ip-store/backend/internal/auth"
	"ip-store/backend/internal/database"
	"ip-store/backend/internal/model"
)

var ( // Define common errors
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("user with this email already exists")
)

type UserService struct {
	db *database.DBContext
}

func NewUserService(db *database.DBContext) *UserService {
	return &UserService{db: db}
}

// CreateUser creates a new user in the database.
func (s *UserService) CreateUser(ctx context.Context, email, password string) (*model.User, error) {
	// Check if user already exists
	_, err := s.GetUserByEmail(ctx, email)
	if err == nil {
		return nil, ErrUserExists
	}
	// If error is not ErrUserNotFound, then it's another DB error
	if err != nil && err != ErrUserNotFound {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
	}

	// Insert user into database
	result, err := s.db.ExecContext(
		ctx,
		"INSERT INTO users (email, password_hash, created_at) VALUES (?, ?, ?)",
		user.Email,
		user.PasswordHash,
		user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	user.ID = lastInsertID

	return user, nil
}

// GetUserByEmail retrieves a user by their email address.
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}
	err := s.db.QueryRowContext(
		ctx,
		"SELECT id, email, password_hash, created_at FROM users WHERE email = ?",
		email,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

// ComparePassword compares a plaintext password with a hashed password.
func (s *UserService) ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (s *UserService) LoginUser(ctx context.Context, email, password string) (string, error) {
	user, err := s.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	err = s.ComparePassword(user.PasswordHash, password)
	if err != nil {
		return "", err
	}

	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}
