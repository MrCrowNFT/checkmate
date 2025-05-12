package service

import (
	"checkmate/api/internal/model"
	"checkmate/api/internal/storage"
	"context"
	"fmt"
)

// get all platform credentials -> for when loading the dashoard
func GetPlatformCredentials(ctx context.Context, userID string) ([]model.PlatformCredential, error) {
	query := `SELECT id, user_id, platform, api_key, created_at 
			FROM platform_credentials
			WHERE user_id = ?;`

	rows, err := storage.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query platform credentials: %w", err)
	}
	defer rows.Close()

	//this what we'll be returning -> array of all credentials
	var credentials []model.PlatformCredential

	//copy the result of row in cred and append it to credentials the got the next row
	for rows.Next() {
		var cred model.PlatformCredential
		if err := rows.Scan(&cred.ID, &cred.UserID, &cred.Platform, &cred.Name, &cred.APIKey, &cred.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan credential row: %w", err)
		}
		credentials = append(credentials, cred)
	}

	return credentials, nil

}
