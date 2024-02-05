package services

import (
	"database/sql"
	"errors"
	"fmt"
	"job-portal-api/internal/models"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// UserService handles business logic related to user operations.
type UserService struct {
	db *sql.DB
}

// NewUserService creates a new UserService instance.
func NewUserService(db *sql.DB) (*UserService, error) {
	// Check if the database connection is nil
	if db == nil {
		return nil, errors.New("db connection cannot be nil")
	}
	return &UserService{db: db}, nil
}

// Create generates a new user record in the database.
func (us *UserService) Create(email, password, role string) (*models.User, error) {
	// Convert email and role to lowercase
	email = strings.ToLower(email)
	role = strings.ToLower(role)

	// Hash the user's password using bcrypt
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	passwordHash := string(hashedBytes)

	// Create a new user instance with hashed password and role
	user := models.User{
		Email:        email,
		PasswordHash: passwordHash,
		Role:         role,
	}

	// Execute the SQL query to insert a new user and retrieve the generated ID
	row := us.db.QueryRow(`
		INSERT INTO users (email, password_hash, role)
		VALUES ($1, $2, $3) RETURNING id`, email, passwordHash, role)
	err = row.Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return &user, nil
}

// Authenticate verifies user credentials and returns the user if authentication is successful.
func (us *UserService) Authenticate(email, password, role string) (*models.User, error) {
	// Convert email and role to lowercase
	email = strings.ToLower(email)
	role = strings.ToLower(role)

	// Create a user instance to store authentication details
	user := models.User{
		Email: email,
		Role:  role,
	}

	// Execute the SQL query to retrieve user information by email
	row := us.db.QueryRow(`
	SELECT id, password_hash, role
	FROM users WHERE email=$1`, email)
	err := row.Scan(&user.ID, &user.PasswordHash, &user.Role)
	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}

	// Check if the provided role matches the user's role
	if user.Role != role {
		return nil, fmt.Errorf("authenticate: invalid role")
	}

	// Compare the provided password with the hashed password in the database
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}

	// Authentication successful, return the user
	return &user, nil
}
