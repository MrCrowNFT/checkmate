package storage

import (
	"database/sql"
	"log"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDb() {
	dbPath := filepath.Join(".", "checkmate.db")

	DB, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("Failed to open db:", err)
	}

	if err := DB.Ping(); err != nil {
		log.Fatal("Failed to connect to db:", err)
	}

	createTables()
}

// creates tables only if the don't exist
func createTables() {
	//don't need to save password, will use firebase auth to handle it
	_, err := DB.Exec(
		//index the email for faster lookup
		`
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY, 
			email TEXT NOT NULL UNIQUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			INDEX idx_email (email)
		);

		CREATE TABLE IF NOT EXISTS tokens (
			id INTEGER PRIMARY KEY,
			user_id TEXT,
			provider TEXT,
			token TEXT,
			FOREIGN KEY(user_id) REFERENCES users(id)
		);
		`)
	if err != nil {
		log.Fatal("Failed to create tables:", err)
	}
}
