package service

import (
	"checkmate/api/internal/model"
	"checkmate/api/internal/storage"
	"context"
	"time"
)

//User CRUD operations

// get user byt id from db
func GetUserById(ctx context.Context, id string) (*model.User, error) {

	var user model.User
	query := `SELECT * FROM users WHERE id= ?`

	err := storage.DB.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.DisplayName,
		&user.CreatedAt,
	)

	if err != nil {
		// In case user not found
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// create new user in db
func CreateUser(ctx context.Context, user *model.User) error {
	db := storage.DB

	//creation time if not already set
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}

	// UPSERT pattern to handle insert and update
	query := `
		INSERT INTO users (id, email, display_name, created_at)
		VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			email = VALUES(email),
			display_name = VALUES(display_name)
	`

	_, err := db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.DisplayName,
		user.CreatedAt,
	)

	return err

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
