package service

import (
	"checkmate/api/internal/model"
	"checkmate/api/internal/platform"
	"checkmate/api/internal/storage"
	"checkmate/api/internal/utils"
	"context"
	"fmt"
	"time"
)

//Credentials CRUD operations

// get all platform credentials -> for when loading the dashoard
func GetPlatformCredentials(ctx context.Context, userID string) ([]model.PlatformCredential, error) {
	query := `SELECT id, user_id, platform, name, api_key, created_at 
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
		//decrypt the api key
		//this should only be used internally so there should not be a problem
		cred.APIKey, err = utils.DecryptString(cred.APIKey)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt API key: %w", err)
		}
		credentials = append(credentials, cred)
	}

	return credentials, nil

}

// not sure if really need this one just yet
func GetPlatformCredentialByID(ctx context.Context, id int, userID string) (*model.PlatformCredential, error) {
	query := `SELECT id, user_id, platform, name, api_key, created_at 
              FROM platform_credentials
              WHERE id = ? AND user_id = ?;`

	var cred model.PlatformCredential
	err := storage.DB.QueryRowContext(ctx, query, id, userID).Scan(
		&cred.ID, &cred.UserID, &cred.Platform, &cred.Name, &cred.APIKey, &cred.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to get credential: %w", err)
	}

	cred.APIKey, err = utils.DecryptString(cred.APIKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt API key: %w", err)
	}

	return &cred, nil
}

// create new plstform credentials
func CreatePlatformCredential(ctx context.Context, userID string, input *model.PlatformCredentialInput) (*model.PlatformCredential, error) {
	//validate credentials before creating cred
	err := validateCredential(ctx, input.Platform, input.APIKey)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials: %w", err)
	}

	query := `INSERT INTO platform_credentials (user_id, platform, name, api_key, created_at)
          VALUES (?, ?, ?, ?, ?)`

	now := time.Now()

	//encrypt the api key before storage
	encryptedAPIKey, err := utils.EncryptString(input.APIKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt API key: %w", err)
	}

	result, err := storage.DB.ExecContext(
		ctx, query, userID, input.Platform, input.Name, encryptedAPIKey, now)
	if err != nil {
		return nil, fmt.Errorf("failed to create platform credential: %w", err)
	}

	//get the id of the last insert to return it to the user
	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	return &model.PlatformCredential{
		ID:        int(id),
		UserID:    userID,
		Platform:  input.Platform,
		Name:      input.Name,
		APIKey:    encryptedAPIKey,
		CreatedAt: now,
	}, nil

}

// for updating platform credential ->needs cred id and user id and the whole PlatformCredentialInput (keep this in mind for the frontend)
func UpdatePlatformCredential(ctx context.Context, id int, userID string, input *model.PlatformCredentialInput) error {
	query := `UPDATE platform_credentials
              SET platform = ?, name = ?, api_key = ?
              WHERE id = ? AND user_id = ?`

	result, err := storage.DB.ExecContext(
		ctx, query, input.Platform, input.Name, input.APIKey, id, userID)

	if err != nil {
		return fmt.Errorf("failed to update platform credential: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("credential not found or you don't have permission to update it")
	}

	return nil
}

// delete a platform credential
func DeletePlatformCredential(ctx context.Context, id int, userID string) error {
	query := `DELETE FROM platform_credentials WHERE id = ? AND user_id = ?`

	result, err := storage.DB.ExecContext(ctx, query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete platform credential: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("credential not found or you don't have permission to delete it")
	}

	return nil
}

func validateCredential(ctx context.Context, platformName string, apiKey string) error {
	switch platformName {
	case "render":
		client := platform.NewRenderClient(apiKey)
		err := client.VerifyCredentials(ctx)
		return err
	case "vercel":
		// TODO: Implement Vercel validation
		return fmt.Errorf("vercel validation not implemented")
	default:
		return fmt.Errorf("unsupported platform: %s", platformName)
	}
}
