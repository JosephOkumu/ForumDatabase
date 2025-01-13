package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func AddReactionHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var data struct {
			UserID       int    `json:"user_id"`
			PostID       *int   `json:"post_id,omitempty"`
			CommentID    *int   `json:"comment_id,omitempty"`
			ReactionType string `json:"reaction_type"` // "LIKE" or "DISLIKE"
		}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}

		if data.UserID == 0 || data.ReactionType == "" ||
			(data.PostID == nil && data.CommentID == nil) {
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}

		if data.ReactionType != "LIKE" && data.ReactionType != "DISLIKE" {
			http.Error(w, "Invalid reaction type", http.StatusBadRequest)
			return
		}

		var table, column string
		var id interface{}

		if data.PostID != nil {
			table = "post_reactions"
			column = "post_id"
			id = data.PostID
		} else if data.CommentID != nil {
			table = "comment_reactions"
			column = "comment_id"
			id = data.CommentID
		}

		// Remove existing reaction for the user on the same post/comment
		deleteQuery := `DELETE FROM ` + table + ` WHERE user_id = ? AND ` + column + ` = ?`
		_, err := db.Exec(deleteQuery, data.UserID, id)
		if err != nil {
			log.Printf("Failed to remove existing reaction: %v\n", err)
			http.Error(w, "Failed to process reaction", http.StatusInternalServerError)
			return
		}

		// Insert the new reaction
		insertQuery := `INSERT INTO ` + table + ` (user_id, ` + column + `, reaction_type) VALUES (?, ?, ?)`
		_, err = db.Exec(insertQuery, data.UserID, id, data.ReactionType)
		if err != nil {
			log.Printf("Failed to add reaction: %v\n", err)
			http.Error(w, "Failed to process reaction", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Reaction added successfully"))
	}
}


