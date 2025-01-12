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
