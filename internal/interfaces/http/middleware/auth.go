package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

func AuthMiddleware(jwtKey []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip auth for public endpoints
			if isPublicPath(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			claims := token.Claims.(jwt.MapClaims)
			userID := claims["user_id"].(string)

			// Add user ID to request context
			ctx := context.WithValue(r.Context(), "user_id", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func isPublicPath(path string) bool {
	publicPaths := []string{
		"/api/auth/register",
		"/api/auth/login",
		"/",
		"/static/",
	}

	for _, pp := range publicPaths {
		if strings.HasPrefix(path, pp) {
			return true
		}
	}
	return false
}
