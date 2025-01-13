package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func AddCommentReactionHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var data struct {
			UserID       int    `json:"user_id"`
			CommentID    int    `json:"comment_id"`
			ReactionType string `json:"reaction_type"` // "LIKE" or "DISLIKE"
		}

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}

		if data.UserID == 0 || data.CommentID == 0 || data.ReactionType == "" {
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}

		if data.ReactionType != "LIKE" && data.ReactionType != "DISLIKE" {
			http.Error(w, "Invalid reaction type", http.StatusBadRequest)
			return
		}

		// Remove existing reaction for the user on the same comment
		deleteQuery := `DELETE FROM comment_reactions WHERE user_id = ? AND comment_id = ?`
		_, err := db.Exec(deleteQuery, data.UserID, data.CommentID)
		if err != nil {
			log.Printf("Failed to remove existing reaction: %v\n", err)
			http.Error(w, "Failed to process reaction", http.StatusInternalServerError)
			return
		}

		// Insert the new reaction
		insertQuery := `INSERT INTO comment_reactions (user_id, comment_id, reaction_type) VALUES (?, ?, ?)`
		_, err = db.Exec(insertQuery, data.UserID, data.CommentID, data.ReactionType)
		if err != nil {
			log.Printf("Failed to add reaction: %v\n", err)
			http.Error(w, "Failed to process reaction", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Reaction added successfully"))
	}
}


