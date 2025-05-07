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
	var err error

	DB, err = sql.Open("sqlite3", dbPath)
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
		//todo: index the email for faster lookup->INDEX idx_email (email)
		//todo: can't cirrently since i am using sqlite 
		`
		CREATE TABLE IF NOT EXISTS users (
			id VARCHAR(128) PRIMARY KEY, 
          	email VARCHAR(255) NOT NULL,
          	display_name VARCHAR(255),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS tokens (
			id VARCHAR(128) PRIMARY KEY,
			user_id VARCHAR(128),
			provider TEXT,
			token TEXT,
			FOREIGN KEY(user_id) REFERENCES users(id)
		);
		`)
	if err != nil {
		log.Fatal("Failed to create tables:", err)
	}
}
