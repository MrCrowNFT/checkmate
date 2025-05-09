package service

import (
	"checkmate/api/internal/model"
	"checkmate/api/internal/storage"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

//User CRUD operations

// get user byt id from db
func GetUserById(ctx context.Context, id string) (*model.User, error) {

	query := `SELECT id, email, display_name, created_at FROM users WHERE id = ?`

	var user model.User
	err := storage.DB.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.DisplayName,
		&user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		// not found but not an error
		log.Printf("User not found with ID: %s", id)
		return nil, nil
	} else if err != nil {
		// database error
		log.Printf("Error querying user: %v", err)
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &user, nil
}

// create new user in db
func CreateUser(ctx context.Context, user *model.User) error {
	log.Printf("Creating new user with ID: %s, Email: %s", user.ID, user.Email)

	// creation time if not set
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}

	query := `INSERT INTO users (id, email, display_name, created_at) VALUES (?, ?, ?, ?)`

	_, err := storage.DB.ExecContext(
		ctx,
		query,
		user.ID,
		user.Email,
		user.DisplayName,
		user.CreatedAt,
	)

	if err != nil {
		log.Printf("Error creating user: %v", err)
		return fmt.Errorf("failed to create user: %w", err)
	}

	log.Printf("Successfully created user with ID: %s", user.ID)
	return nil
}

// update user in db
func UpdateUser(ctx context.Context, user *model.User) error {
	db := storage.DB

	query := `
		UPDATE users
		SET email = ?, display_name = ?
		WHERE id = ?
	`

	_, err := db.ExecContext(ctx, query,
		user.Email,
		user.DisplayName,
		user.ID,
	)

	return err
}

func DeleteUser(ctx context.Context, id string) error {
	db := storage.DB

	query := `DELETE FROM users WHERE id = ?`

	_, err := db.ExecContext(ctx, query, id)

	return err
}
