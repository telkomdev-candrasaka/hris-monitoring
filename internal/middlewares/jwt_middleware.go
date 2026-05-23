package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const ContextUserKey contextKey = "jwtUser"

type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Authorization header required"})
			return
		}

		parts := strings.Fields(authHeader)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid authorization header format"})
			return
		}

		tokenString := parts[1]
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "supersecret"
		}

		claims := &JWTClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid or expired token"})
			return
		}

		ctx := context.WithValue(r.Context(), ContextUserKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RoleAuthorize(allowedRoles ...string) func(http.Handler) http.Handler {
	roles := map[string]struct{}{}
	for _, role := range allowedRoles {
		roles[role] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			value := r.Context().Value(ContextUserKey)
			claims, ok := value.(*JWTClaims)
			if !ok || claims == nil {
				writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Token tidak ditemukan atau tidak valid"})
				return
			}

			if _, allowed := roles[claims.Role]; !allowed {
				writeJSON(w, http.StatusForbidden, map[string]string{"error": "Akses ditolak"})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}
