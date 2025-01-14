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
		var hashedPassword, username string
		query := "SELECT id, password, username FROM users WHERE email = ?"
		err = db.QueryRow(query, req.Email).Scan(&userID, &hashedPassword, &username)
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

		// Respond with a success message and the username
		response := struct {
			Message  string `json:"message"`
			Username string `json:"username"`
		}{
			Message:  "Login successful",
			Username: username,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}


// LogoutHandler handles user logout requests
func LogoutHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the session token from the cookie
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, "Session not found", http.StatusUnauthorized)
			return
		}

		// Delete session from the database
		_, err = db.Exec("DELETE FROM sessions WHERE uuid = ?", cookie.Value)
		if err != nil {
			http.Error(w, "Failed to logout", http.StatusInternalServerError)
			return
		}

		// Remove the session cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    "",
			Expires:  time.Unix(0, 0), // Expire the cookie
			HttpOnly: true,
		})

		// Respond with a success message
		response := struct {
			Message string `json:"message"`
		}{
			Message: "Logout successful",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
