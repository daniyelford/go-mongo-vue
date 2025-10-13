package middleware

import (
	"context"
	"go-mongo-vue-go/config"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const MobileContextKey contextKey = "mobile"

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return config.JwtSecret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		mobile, ok := claims["mobile"].(string)
		if !ok || mobile == "" {
			http.Error(w, "invalid token claims", http.StatusUnauthorized)
			return
		}
		val, err := config.RedisClient.Get(r.Context(), "token:"+mobile).Result()
		if err != nil || val != tokenString {
			http.Error(w, "expired or revoked token", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), MobileContextKey, mobile)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
