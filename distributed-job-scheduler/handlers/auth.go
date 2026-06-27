package handlers

import (
	"encoding/json"
	"net/http"
	"errors"
	"time"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"


)

func RegisterHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if input.Email == "" || input.Password == "" {
			http.Error(w, "email and password are required", http.StatusBadRequest)
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "failed to hash password", http.StatusInternalServerError)
			return
		}

		var id int64
		err = pool.QueryRow(
			r.Context(),
			"INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id",
			input.Email, string(hash),
		).Scan(&id)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				http.Error(w, "email already exists", http.StatusConflict)
				return
			}
			http.Error(w, "failed to create user", http.StatusInternalServerError)
			return

		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"id":    id,
			"email": input.Email,
		})
	}
}

func LoginHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		var userID int64
		var passwordHash string
		err = pool.QueryRow(
			r.Context(),
			"SELECT id, password_hash FROM users WHERE email = $1",
			input.Email,
		).Scan(&userID, &passwordHash)

		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "invalid email or password", http.StatusUnauthorized)
			return
		}
		if err != nil {
			http.Error(w, "failed to look up user", http.StatusInternalServerError)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(input.Password))
		if err != nil {
			http.Error(w, "invalid email or password", http.StatusUnauthorized)
			return
		}

		claims := jwt.MapClaims{
			"user_id": userID,
			"exp":     time.Now().Add(24 * time.Hour).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
		if err != nil {
			http.Error(w, "failed to generate token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"token": signedToken,
		})
	}
}