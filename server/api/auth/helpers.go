package authAPI

import (
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func GetUserIDFromRequest(r *http.Request) (uint, bool) {
	id, ok := r.Context().Value(UserIDContextKey).(uint)
	return id, ok
}

func ParseToken(tokenStr string) (uint, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtSecretBytes, nil
	})
	if err != nil || !token.Valid {
		return 0, errors.New("invalid token")
	}
	return claims.UserID, nil
}
