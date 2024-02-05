package middleware

import (
	"context"
	"net/http"

	"github.com/rs/zerolog/log"
)

// JWTMiddlewareCookie is a middleware function that checks for a JWT token in the request cookie,
// verifies the token, and sets the user ID in the request context if the token is valid.
// It also performs role-based authorization by checking the required role against the token claims.
func (m Mid) JWTMiddlewareCookie(next http.HandlerFunc, requiredRole string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the JWT token from the request cookie
		cookie, err := r.Cookie("token")
		if err != nil {
			log.Error().Err(err).Send()
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// Verify the JWT token using the authentication service
		claim, err := m.a.VerifyToken(cookie.Value, requiredRole)
		if err != nil {
			log.Error().Err(err).Send()
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// You can use the claims for further authorization checks

		// Set the user ID from the token claims in the request context
		ctx := context.WithValue(r.Context(), "userID", claim.Subject)

		// Call the next handler in the chain with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
