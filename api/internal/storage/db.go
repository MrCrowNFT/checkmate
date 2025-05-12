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
func createTables() error {
	//don't need to save password, will use firebase auth to handle it
	_, err := DB.Exec(
		//todo: index the email for faster lookup->INDEX idx_email (email)
		//todo: can't currently since i am using sqlite
		`
		-- Users table;
	CREATE TABLE IF NOT EXISTS users (
    	id VARCHAR(128) PRIMARY KEY, 
    	email VARCHAR(255) NOT NULL,
	    display_name VARCHAR(255),
    	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- Platform credentials;
	CREATE TABLE IF NOT EXISTS platform_credentials (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
	    user_id VARCHAR(128) NOT NULL,
    	platform VARCHAR(50) NOT NULL,  
	    name VARCHAR(255) NOT NULL,     
	    api_key TEXT NOT NULL,          
	    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	-- Cache for deployment data
	CREATE TABLE IF NOT EXISTS deployment_cache (
    	id VARCHAR(255) NOT NULL,       
    	platform_credential_id INTEGER NOT NULL,
    	name VARCHAR(255) NOT NULL,
    	status VARCHAR(50) NOT NULL,    
    	url VARCHAR(255),
    	last_deployed_at TIMESTAMP,
	    branch VARCHAR(255),
    	service_type VARCHAR(100),      
    	framework VARCHAR(100),         
	    last_updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    	metadata TEXT,                  
	    PRIMARY KEY (id, platform_credential_id),
    	FOREIGN KEY (platform_credential_id) REFERENCES platform_credentials(id) ON DELETE CASCADE
	);
	
		`)
	if err != nil {
		log.Printf("Failed to create tables: %v", err)
		return err
	}

	return nil
}
