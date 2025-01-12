package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// DB is a global variable for database connection
var DB *sql.DB

// Initialize initializes the database connection and applies the schema
func Initialize() error {
	var err error
	DB, err = sql.Open("sqlite3", "./forum.db")
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	// Ensure the database is accessible
	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	// Apply the schema
	err = applySchema()
	if err != nil {
		return fmt.Errorf("failed to apply schema: %v", err)
	}

	log.Println("Database initialized successfully")
	return nil
}

// applySchema applies the SQL schema from the schema.sql file
func applySchema() error {
	schemaContent, err := os.ReadFile("./db/schema.sql")
	if err != nil {
		return fmt.Errorf("failed to read schema file: %v", err)
	}

	_, err = DB.Exec(string(schemaContent))
	if err != nil {
		return fmt.Errorf("failed to execute schema SQL: %v", err)
	}

	log.Println("Schema applied successfully")
	return nil
}


