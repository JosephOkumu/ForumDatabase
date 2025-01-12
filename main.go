package main

import (
	"fmt"
	"log"
	"net/http"

	"forum/db"
	"forum/handlers"
)

func main() {
	// Initialize the database
	err := db.Initialize()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer db.Close() // Ensure database connection is closed

	
	// Define routes
	http.HandleFunc("/register", handlers.RegisterUserHandler) 
    http.HandleFunc("/login", handlers.LoginHandler(db.DB))
	http.HandleFunc("/posts", handlers.GetPostsHandler(db.DB))
	http.HandleFunc("/create-post", handlers.CreatePostHandler(db.DB))

	// Start the server
	fmt.Println("Server started on :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
