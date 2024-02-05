package middleware

import (
	"errors"
	"job-portal-api/internal/auth"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
)

// Mid is a middleware struct containing an authentication instance.
type Mid struct {
	a *auth.Auth
}

// NewMid creates a new middleware instance with the provided authentication service.
// It returns an error if the authentication service is nil.
func NewMid(a *auth.Auth) (*Mid, error) {
	if a == nil {
		return nil, errors.New("auth struct cannot be nil")
	}
	return &Mid{a: a}, nil
}

// JWTMiddleware is a middleware function that checks for a JWT token in the request header,
// verifies the token, and allows the request to proceed if the token is valid and the required role is satisfied.
func (m Mid) JWTMiddleware(next http.HandlerFunc, requiredRole string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the JWT token from the request header
		tokenString := extractTokenFromHeader(r)
		if tokenString == "" {
			// If no token is found, respond with Unauthorized status
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// Verify the JWT token using the authentication service
		_, err := m.a.VerifyToken(tokenString, requiredRole)
		if err != nil {
			log.Error().Err(err).Send()
			// If token verification fails, respond with Unauthorized status
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

// extractTokenFromHeader extracts the JWT token from the Authorization header in the request.
func extractTokenFromHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}
	// Split the Authorization header into parts
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}
	// Return the token part
	return parts[1]
}
