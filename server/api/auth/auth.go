package authAPI

import (
	"OnlineGame/database"
	APIhelpers "OnlineGame/server/api/helpers"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func RegisterAuthRoutes(r *mux.Router) {
	r.HandleFunc("/auth/register", Register).Methods("POST")
	r.HandleFunc("/auth/login", Login).Methods("POST")
}

func Register(w http.ResponseWriter, r *http.Request) {
	var body struct {
		UserName string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		APIhelpers.ErrorJSONResponse(w, "invalid body", http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		APIhelpers.ErrorJSONResponse(w, "failed to hash password", http.StatusInternalServerError)
		return
	}

	repo := database.NewUserRepository(database.GetDB())
	// Check for unique username
	found, err := repo.FindByUserName(body.UserName)
	if err != nil || found != nil {
		APIhelpers.ErrorJSONResponse(w, "this username already taken", http.StatusNotFound)
		return
	}
	// Check for unique email
	found, err = repo.FindByEmail(body.Email)
	if err != nil || found != nil {
		APIhelpers.ErrorJSONResponse(w, "this email already registered", http.StatusNotFound)
		return
	}
	// Create
	created, err := repo.Create(body.UserName, body.Email, string(hash))
	if err != nil {
		APIhelpers.ErrorJSONResponse(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	token, err := GenerateToken(created.ID)
	if err != nil {
		APIhelpers.ErrorJSONResponse(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	APIhelpers.JSONResponse(w, map[string]string{"token": token}, http.StatusCreated)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		APIhelpers.ErrorJSONResponse(w, "invalid body", http.StatusBadRequest)
		return
	}

	repo := database.NewUserRepository(database.GetDB())
	found, err := repo.FindByEmail(body.Email)
	if err != nil || found == nil {
		APIhelpers.ErrorJSONResponse(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(found.Password), []byte(body.Password)); err != nil {
		APIhelpers.ErrorJSONResponse(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := GenerateToken(found.ID)
	if err != nil {
		APIhelpers.ErrorJSONResponse(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	APIhelpers.JSONResponse(w, map[string]string{"token": token}, http.StatusOK)
}

func GenerateToken(userID uint) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
