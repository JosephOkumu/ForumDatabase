package handlers

import (
	"encoding/json"
	"fmt"
	"forum/db"
	"forum/models"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt" 
)

// InsertUser inserts a new user into the database
func InsertUser(user models.User) error {
	// Check if the username already exists
	var count int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", user.Username).Scan(&count)
	if err != nil {
		log.Printf("Error checking username existence: %v", err)
		return fmt.Errorf("failed to check username existence: %v", err)
	}
	if count > 0 {
		return fmt.Errorf("username already exists")
	}

	// Check if the email already exists
	err = db.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", user.Email).Scan(&count)
	if err != nil {
		log.Printf("Error checking email existence: %v", err)
		return fmt.Errorf("failed to check email existence: %v", err)
	}
	if count > 0 {
		return fmt.Errorf("email already exists")
	}

	// Hash the password before inserting
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return fmt.Errorf("failed to hash password: %v", err)
	}

	// Insert user into the database
	stmt, err := db.DB.Prepare("INSERT INTO users (email, username, password, created_at) VALUES (?, ?, ?, datetime('now'))")
	if err != nil {
		log.Printf("Error preparing SQL statement: %v", err)
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Email, user.Username, string(hashedPassword))
	if err != nil {
		log.Printf("Error executing SQL statement: %v", err)
		return fmt.Errorf("failed to insert user: %v", err)
	}

	log.Println("User inserted successfully")
	return nil
}

// RegisterUserHandler handles user registration
func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var user models.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			log.Printf("Error decoding request body: %v", err)
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		// Insert user into the database
		err = InsertUser(user)
		if err != nil {
			log.Printf("Error inserting user: %v", err)
			http.Error(w, fmt.Sprintf("Error registering user: %v", err), http.StatusInternalServerError)
			return
		}

		// Respond with success
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "User registered successfully")
		return
	}

	http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
}
