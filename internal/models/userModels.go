package models

// NewUser represents the structure for creating a new user. It includes fields such as email, password, and role.
type NewUser struct {
	Email    string `json:"email" validate:"required,email"`              // Email is the email address of the user and is required.
	Password string `json:"password" validate:"required,min=6"`           // Password is the user's password and is required with a minimum length of 6 characters.
	Role     string `json:"role" validate:"required,customRoleValidator"` // Role represents the user's role and is required, with custom validation using the "customRoleValidator" tag.
}

// User represents the structure for a user entity. It includes fields such as ID, email, password hash, and role.
type User struct {
	ID           int    `json:"id"`    // ID is a unique identifier for the user.
	Email        string `json:"email"` // Email is the email address of the user.
	PasswordHash string `json:"-"`     // PasswordHash is the hashed version of the user's password and is not included in JSON responses.
	Role         string `json:"role"`  // Role represents the user's role.
}
