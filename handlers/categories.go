package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func GetPostsByCategoryHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Ensure the method is GET
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Extract category_id from query parameters
		categoryIDStr := r.URL.Query().Get("category_id")
		if categoryIDStr == "" {
			http.Error(w, "Missing category_id parameter", http.StatusBadRequest)
			return
		}

		// Convert category_id to an integer
		categoryID, err := strconv.Atoi(categoryIDStr)
		if err != nil {
			http.Error(w, "Invalid category_id value", http.StatusBadRequest)
			return
		}

		// Query to fetch posts for the given category
		query := `
            SELECT p.id, p.title, p.content, p.created_at, p.author_id
            FROM posts p
            INNER JOIN post_categories pc ON p.id = pc.post_id
            WHERE pc.category_id = ?`

		rows, err := db.Query(query, categoryID)
		if err != nil {
			log.Printf("Failed to fetch posts by category: %v\n", err)
			http.Error(w, "Failed to fetch posts by category", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Parse rows into a slice of posts
		type Post struct {
			ID        int    `json:"id"`
			Title     string `json:"title"`
			Content   string `json:"content"`
			CreatedAt string `json:"created_at"`
			AuthorID  int    `json:"author_id"`
		}
		var posts []Post

		for rows.Next() {
			var post Post
			if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.AuthorID); err != nil {
				log.Printf("Failed to scan post row: %v\n", err)
				http.Error(w, "Failed to process posts", http.StatusInternalServerError)
				return
			}
			posts = append(posts, post)
		}

		// Send posts as JSON response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)
	}
}
