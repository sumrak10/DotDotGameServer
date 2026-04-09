package authAPI

import (
	"OnlineGame/database"
	APIhelpers "OnlineGame/server/api/helpers"
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDContextKey contextKey = "user_id"

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			APIhelpers.ErrorJSONResponse(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			APIhelpers.ErrorJSONResponse(w, "invalid token", http.StatusUnauthorized)
			return
		}

		repo := database.NewUserRepository(database.GetDB())
		foundUser, err := repo.FindByID(claims.UserID)
		if err != nil {
			APIhelpers.ErrorJSONResponse(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		if foundUser == nil {
			APIhelpers.ErrorJSONResponse(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDContextKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
