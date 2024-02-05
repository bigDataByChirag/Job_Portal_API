package handlers

import (
	"encoding/json"
	"errors"
	"job-portal-api/internal/auth"
	"job-portal-api/internal/models"
	"job-portal-api/internal/services"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

// Users struct represents the handler for user-related operations
type Users struct {
	userService *services.UserService
	a           *auth.Auth
}

// NewUsers creates a new Users handler with the provided services and authentication
func NewUsers(us *services.UserService, a *auth.Auth) (*Users, error) {
	if us == nil || a == nil {
		return nil, errors.New("please provide all the values")
	}
	return &Users{
		userService: us,
		a:           a,
	}, nil
}

// customRoleValidator is a custom validation function for user roles
func customRoleValidator(fl validator.FieldLevel) bool {
	role := fl.Field().String()
	return role == "admin" || role == "user"
}

// HandlerError struct represents an error response structure
type HandlerError struct {
	Message string `json:"message"`
}

// CreateUser handles the creation of a new user
func (u Users) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newUser models.NewUser

	// Create a new validator and register the custom role validator
	validate := validator.New()
	validate.RegisterValidation("customRoleValidator", customRoleValidator)

	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		log.Error().Err(err).Send()
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the user struct
	if err := validate.Struct(newUser); err != nil {
		log.Error().Err(err).Send()
		appErr := HandlerError{
			Message: http.StatusText(http.StatusInternalServerError),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(appErr)
		return
	}

	email := newUser.Email
	password := newUser.Password
	role := newUser.Role

	// Create the user using the user service
	_, err = u.userService.Create(email, password, role)
	if err != nil {
		log.Error().Err(err).Send()
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("User Registerd Successfully")
}

// ProcessLoginIn handles the user login process
func (u Users) ProcessLoginIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var authUser struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
		Role     string `json:"role" validate:"required,customRoleValidator"`
	}

	// Create a new validator and register the custom role validator
	validate := validator.New()
	validate.RegisterValidation("customRoleValidator", customRoleValidator)

	err := json.NewDecoder(r.Body).Decode(&authUser)
	if err != nil {
		log.Error().Err(err).Send()
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the user struct
	if err := validate.Struct(authUser); err != nil {
		log.Error().Err(err).Send()
		appErr := HandlerError{
			Message: http.StatusText(http.StatusInternalServerError),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(appErr)
		return
	}

	// Authenticate the user using the user service
	user, err := u.userService.Authenticate(authUser.Email, authUser.Password, authUser.Role)
	if err != nil {
		log.Error().Err(err).Send()
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	// Generate a JWT token for the authenticated user
	tkn, err := u.a.GenerateToken(user.ID, user.Role)
	if err != nil {
		log.Error().Err(err).Send()
		http.Error(w, "Failed to generate JWT token", http.StatusInternalServerError)
		return
	}

	// Set the JWT token as an HTTP cookie
	cookie := http.Cookie{
		Name:     "token",
		Value:    tkn,
		HttpOnly: true, // This ensures that only the browser can access the token, not external scripts
	}
	http.SetCookie(w, &cookie)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("User Logged-In Successfully")
}
