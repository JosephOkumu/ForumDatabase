package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Message string `json:"message"`
}

// LoginHandler handles user login requests
func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Ensure the request method is POST
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse and decode the request body
		var req LoginRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Validate input
		if req.Email == "" || req.Password == "" {
			http.Error(w, "Email and password are required", http.StatusBadRequest)
			return
		}

		// Query the database for the user
		var userID int
		var hashedPassword string
		query := "SELECT id, password FROM users WHERE email = ?"
		err = db.QueryRow(query, req.Email).Scan(&userID, &hashedPassword)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Invalid email or password", http.StatusUnauthorized)
				return
			}
			log.Println("Database query error:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Compare the stored hash with the entered password
		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password))
		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		// Generate a session token
		sessionToken := uuid.New().String()
		expiresAt := time.Now().Add(24 * time.Hour)

		// Insert session into the database
		_, err = db.Exec("INSERT INTO sessions (uuid, user_id, expires_at) VALUES (?, ?, ?)",
			sessionToken, userID, expiresAt)
		if err != nil {
			log.Println("Failed to store session:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Set the session token as a cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    sessionToken,
			Expires:  expiresAt,
			HttpOnly: true, // Prevent JavaScript access for security
		})

		// Respond with a success message
		response := LoginResponse{Message: "Login successful"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
