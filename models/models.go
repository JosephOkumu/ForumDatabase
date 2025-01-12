package models

// User represents a user in the forum
type User struct {
    Email      string `json:"email"`
    Username   string `json:"username"`
    Password   string `json:"password"` // In real apps, use hashed passwords
    
}
