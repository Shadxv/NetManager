package auth

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const userContextKey = contextKey("user")

func AuthenticateToken(jwtSecret string, allowTempUser bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Access denied", http.StatusForbidden)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := ValidateToken(tokenString, jwtSecret)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusForbidden)
				return
			}

			if claims.RequiresPasswordReset && !allowTempUser {
				http.Error(w, "Password reset required", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Authorize(requiredBits ...PermissionFlags) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(userContextKey).(*TokenPayload)
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			userPerms := PermissionFlags(claims.RolePermissions)

			if (userPerms & ADMIN) == ADMIN {
				next.ServeHTTP(w, r)
				return
			}

			for _, bit := range requiredBits {
				if (userPerms & bit) == bit {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, "Forbidden: Insufficient permissions", http.StatusForbidden)
		})
	}
}
