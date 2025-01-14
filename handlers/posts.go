package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// Post represents the structure of a post
type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
}

// CreatePostHandler handles creating a new post
func CreatePostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Incoming request to create post")

		// Ensure the request method is POST
		if r.Method != http.MethodPost {
			log.Println("Invalid method:", r.Method)
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}
		

		// Retrieve session token from cookie
		cookie, err := r.Cookie("session_token")
		if err != nil {
			log.Println("Failed to retrieve session token:", err)
			http.Error(w, "Unauthorized: Please log in first", http.StatusUnauthorized)
			return
		}
		log.Println("Session token received:", cookie.Value)

		// Validate the session token in the database
		var userID int
		query := "SELECT user_id FROM sessions WHERE uuid = ? AND expires_at > DATETIME('now')"
		err = db.QueryRow(query, cookie.Value).Scan(&userID)
		if err != nil {
			log.Println("Session validation failed. Token might be invalid or expired:", err)
			http.Error(w, "Unauthorized: Please log in first", http.StatusUnauthorized)
			return
		}
		log.Println("Session validated successfully. User ID:", userID)

		// Parse and decode the request body
		var req struct {
			Title   string `json:"title"`
			Content string `json:"content"`
		}
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Println("Failed to decode request body:", err)
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		log.Println("Request payload decoded:", req)

		// Validate input
		if req.Title == "" || req.Content == "" {
			log.Println("Validation failed: Title or Content is empty")
			http.Error(w, "Title and content are required", http.StatusBadRequest)
			return
		}

		// Insert the new post into the database
		query = "INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)"
		_, err = db.Exec(query, userID, req.Title, req.Content)
		if err != nil {
			log.Println("Error inserting post into database:", err)
			http.Error(w, "Failed to create post", http.StatusInternalServerError)
			return
		}
		log.Println("Post successfully created for user ID:", userID)

		// Respond with a success message
		w.Header().Set("Content-Type", "application/json")
		response := map[string]string{
			"message": "Post created successfully",
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Println("Failed to encode response:", err)
			http.Error(w, "Failed to create post", http.StatusInternalServerError)
		}
	}
}

// GetPostsHandler handles fetching posts from the database
func GetPostsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Incoming request to fetch posts")

		// Ensure the request method is GET
		if r.Method != http.MethodGet {
			log.Println("Invalid method:", r.Method)
			http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
			return
		}

		// Query to fetch posts with author name
		query := `
			SELECT 
				posts.id, 
				posts.title, 
				posts.content, 
				users.username AS author, 
				posts.created_at
			FROM 
				posts
			JOIN 
				users 
			ON 
				posts.user_id = users.id
			ORDER BY 
				posts.created_at DESC
		`

		// Execute the query
		rows, err := db.Query(query)
		if err != nil {
			log.Println("Error executing query:", err)
			http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Prepare the response
		var posts []Post
		for rows.Next() {
			var post Post
			err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Author, &post.CreatedAt)
			if err != nil {
				log.Println("Error scanning row:", err)
				http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
				return
			}
			posts = append(posts, post)
		}

		// Check for errors after row iteration
		if err = rows.Err(); err != nil {
			log.Println("Row iteration error:", err)
			http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
			return
		}

		// Respond with the posts in JSON format
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(posts); err != nil {
			log.Println("Failed to encode response:", err)
			http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		}
	}
}
