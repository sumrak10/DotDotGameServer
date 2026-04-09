package tests

import (
	"OnlineGame/server/api/auth"
	"testing"
)

func generateTestToken(t *testing.T, userID uint) string {
	token, err := auth.GenerateToken(userID)
	if err != nil {
		t.Fatal(err)
	}
	return token
}
