package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// Comment represents a single comment on a post.
type Comment struct {
	ID        int    `json:"id"`
	PostID    int    `json:"post_id"`
	UserID    int    `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

// AddCommentHandler allows a user to add a comment to a post and immediately returns the new comment.
func AddCommentHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Ensure the method is POST
        if r.Method != http.MethodPost {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }

        // Log method and request URL for debugging
        log.Printf("Received POST request to %s", r.URL.Path)

        // Parse JSON body
        var data struct {
            PostID  int    `json:"post_id"`
            Content string `json:"content"`
        }

        decoder := json.NewDecoder(r.Body)
        err := decoder.Decode(&data)
        if err != nil {
            log.Printf("Failed to decode JSON: %v", err)
            http.Error(w, "Invalid JSON body", http.StatusBadRequest)
            return
        }

        // Log parsed data
        log.Printf("Parsed Data: %+v", data)

        // Validate post_id and content
        if data.PostID == 0 || data.Content == "" {
            log.Println("Missing post_id or content in request")
            http.Error(w, "Missing post_id or content", http.StatusBadRequest)
            return
        }

        // Insert the comment into the database
        query := "INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)"
        res, err := db.Exec(query, data.PostID, 0, data.Content) // Using userID = 0 as a placeholder
        if err != nil {
            log.Printf("Failed to insert comment into database: %v", err)
            http.Error(w, "Failed to add comment", http.StatusInternalServerError)
            return
        }

        // Retrieve the inserted comment ID
        commentID, err := res.LastInsertId()
        if err != nil {
            log.Printf("Failed to retrieve inserted comment ID: %v", err)
            http.Error(w, "Failed to retrieve comment", http.StatusInternalServerError)
            return
        }

        // Fetch the newly inserted comment from the database
        query = "SELECT id, post_id, user_id, content, created_at FROM comments WHERE id = ?"
        var comment Comment
        err = db.QueryRow(query, commentID).Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt)
        if err != nil {
            log.Printf("Failed to retrieve inserted comment: %v", err)
            http.Error(w, "Failed to retrieve comment", http.StatusInternalServerError)
            return
        }

        // Log the retrieved comment
        log.Printf("Successfully inserted comment: %+v", comment)

        // Return the new comment as JSON
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(comment)
    }
}

// GetCommentsHandler retrieves all comments for a specific post.
func GetCommentsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Ensure the method is GET
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse the post ID from query parameters
		postIDStr := r.URL.Query().Get("post_id")
		if postIDStr == "" {
			http.Error(w, "Missing post_id query parameter", http.StatusBadRequest)
			return
		}

		// Convert postID to an integer
		postID, err := strconv.Atoi(postIDStr)
		if err != nil {
			http.Error(w, "Invalid post_id", http.StatusBadRequest)
			return
		}

		// Query the database for comments related to the post
		query := "SELECT id, post_id, user_id, content, created_at FROM comments WHERE post_id = ? ORDER BY created_at ASC"
		rows, err := db.Query(query, postID)
		if err != nil {
			log.Printf("Failed to retrieve comments: %v\n", err)
			http.Error(w, "Failed to retrieve comments", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Parse the results
		var comments []Comment
		for rows.Next() {
			var comment Comment
			if err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt); err != nil {
				log.Printf("Failed to parse comment: %v\n", err)
				http.Error(w, "Failed to retrieve comments", http.StatusInternalServerError)
				return
			}
			comments = append(comments, comment)
		}

		// Return the comments as JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(comments)
	}
}
