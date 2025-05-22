package service

import (
	"checkmate/api/internal/model"
	"checkmate/api/internal/platform"
	"checkmate/api/internal/storage"
	"checkmate/api/internal/utils"
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

//Credentials CRUD operations

// get all platform credentials -> for getting all the deployments -> when loading dashboard
// should only be used internally
func GetPlatformCredentials(ctx context.Context, userID string) ([]model.PlatformCredential, error) {
	logger := log.WithFields(log.Fields{
		"func":       "GetPlatformCredentials",
		"user_id":    userID,
		"request_id": ctx.Value("request_id"),
	})

	logger.Debug("Getting platform credentials started")

	query := `SELECT id, user_id, platform, api_key, created_at 
        FROM platform_credentials
        WHERE user_id = ?;`

	rows, err := storage.DB.QueryContext(ctx, query, userID)
	if err != nil {
		logger.WithError(err).Error("Failed to query platform credentials")
		return nil, fmt.Errorf("failed to query platform credentials: %w", err)
	}
	defer rows.Close()

	//this what we'll be returning -> array of all credentials
	var credentials []model.PlatformCredential

	//copy the result of row in cred and append it to credentials the got the next row
	for rows.Next() {
		var cred model.PlatformCredential
		if err := rows.Scan(&cred.ID, &cred.UserID, &cred.Platform, &cred.APIKey, &cred.CreatedAt); err != nil {
			logger.WithError(err).Error("Failed to scan credential row")
			return nil, fmt.Errorf("failed to scan credential row: %w", err)
		}
		//decrypt the api key
		//this should only be used internally so there should not be a problem
		cred.APIKey, err = utils.DecryptString(cred.APIKey)
		if err != nil {
			logger.WithError(err).Error("Failed to decrypt API key")
			return nil, fmt.Errorf("failed to decrypt API key: %w", err)
		}
		credentials = append(credentials, cred)
	}

	logger.WithField("credentials_count", len(credentials)).Debug("Retrieved platform credentials successfully")
	return credentials, nil
}

// not sure if really need this one just yet -> yeah, i really don´t need this one
func GetPlatformCredentialByID(ctx context.Context, id int, userID string) (*model.PlatformCredential, error) {
	logger := log.WithFields(log.Fields{
		"func":          "GetPlatformCredentialByID",
		"credential_id": id,
		"user_id":       userID,
		"request_id":    ctx.Value("request_id"),
	})

	logger.Debug("Getting platform credential by ID started")

	query := `SELECT id, user_id, platform, api_key, created_at 
              FROM platform_credentials
              WHERE id = ? AND user_id = ?;`

	var cred model.PlatformCredential
	err := storage.DB.QueryRowContext(ctx, query, id, userID).Scan(
		&cred.ID, &cred.UserID, &cred.Platform, &cred.APIKey, &cred.CreatedAt)

	if err != nil {
		logger.WithError(err).Error("Failed to get credential")
		return nil, fmt.Errorf("failed to get credential: %w", err)
	}

	cred.APIKey, err = utils.DecryptString(cred.APIKey)
	if err != nil {
		logger.WithError(err).Error("Failed to decrypt API key")
		return nil, fmt.Errorf("failed to decrypt API key: %w", err)
	}

	logger.WithFields(log.Fields{
		"platform": cred.Platform,
	}).Debug("Retrieved platform credential successfully")
	return &cred, nil
}

// create new plstform credentials
func CreatePlatformCredential(ctx context.Context, userID string, input *model.PlatformCredentialInput) (*model.PlatformCredential, error) {
	logger := log.WithFields(log.Fields{
		"func":       "CreatePlatformCredential",
		"user_id":    userID,
		"platform":   input.Platform,
		"request_id": ctx.Value("request_id"),
	})

	logger.Debug("Creating platform credential started")

	//validate credentials before creating cred
	//todo considere deleting this one, we already validate on the handler
	err := ValidateCredential(ctx, input.Platform, input.APIKey)
	if err != nil {
		logger.WithError(err).Warn("Validation failed for credential")
		return nil, fmt.Errorf("invalid credentials: %w", err)
	}

	logger.Debug("Credential validated successfully")

	query := `INSERT INTO platform_credentials (user_id, platform, api_key, created_at)
          VALUES (?, ?, ?, ?, ?)`

	now := time.Now()

	//encrypt the api key before storage
	encryptedAPIKey, err := utils.EncryptString(input.APIKey)
	if err != nil {
		logger.WithError(err).Error("Failed to encrypt API key")
		return nil, fmt.Errorf("failed to encrypt API key: %w", err)
	}

	logger.Debug("API key encrypted successfully")

	result, err := storage.DB.ExecContext(
		ctx, query, userID, input.Platform, encryptedAPIKey, now)
	if err != nil {
		logger.WithError(err).Error("Failed to create platform credential in database")
		return nil, fmt.Errorf("failed to create platform credential: %w", err)
	}

	//get the id of the last insert to return it to the user
	id, err := result.LastInsertId()
	if err != nil {
		logger.WithError(err).Error("Failed to get last insert ID")
		return nil, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	logger.WithField("credential_id", id).Info("Platform credential created successfully")
	return &model.PlatformCredential{
		ID:        int(id),
		UserID:    userID,
		Platform:  input.Platform,
		APIKey:    encryptedAPIKey,
		CreatedAt: now,
	}, nil
}

// for updating platform credential ->needs cred id and user id and the whole PlatformCredentialInput (keep this in mind for the frontend)
func UpdatePlatformCredential(ctx context.Context, id int, userID string, input *model.PlatformCredentialInput) error {
	logger := log.WithFields(log.Fields{
		"func":          "UpdatePlatformCredential",
		"credential_id": id,
		"user_id":       userID,
		"platform":      input.Platform,
		"request_id":    ctx.Value("request_id"),
	})

	logger.Debug("Updating platform credential started")

	encryptedAPIKey, err := utils.EncryptString(input.APIKey)
	if err != nil {
		logger.WithError(err).Error("Failed to encrypt API key")
		return fmt.Errorf("failed to encrypt API key: %w", err)
	}

	logger.Debug("API key encrypted successfully")

	query := `UPDATE platform_credentials
              SET platform = ?, api_key = ?
              WHERE id = ? AND user_id = ?`

	result, err := storage.DB.ExecContext(
		ctx, query, input.Platform, encryptedAPIKey, id, userID)

	if err != nil {
		logger.WithError(err).Error("Failed to update platform credential in database")
		return fmt.Errorf("failed to update platform credential: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.WithError(err).Error("Failed to get rows affected")
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		logger.Warn("Credential not found or user doesn't have permission to update it")
		return fmt.Errorf("credential not found or you don't have permission to update it")
	}

	logger.Info("Platform credential updated successfully")
	return nil
}

// delete a platform credential
func DeletePlatformCredential(ctx context.Context, id int, userID string) error {
	logger := log.WithFields(log.Fields{
		"func":          "DeletePlatformCredential",
		"credential_id": id,
		"user_id":       userID,
		"request_id":    ctx.Value("request_id"),
	})

	logger.Debug("Deleting platform credential started")

	query := `DELETE FROM platform_credentials WHERE id = ? AND user_id = ?`

	result, err := storage.DB.ExecContext(ctx, query, id, userID)
	if err != nil {
		logger.WithError(err).Error("Failed to delete platform credential from database")
		return fmt.Errorf("failed to delete platform credential: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.WithError(err).Error("Failed to get rows affected")
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		logger.Warn("Credential not found or user doesn't have permission to delete it")
		return fmt.Errorf("credential not found or you don't have permission to delete it")
	}

	logger.Info("Platform credential deleted successfully")
	return nil
}

func ValidateCredential(ctx context.Context, platformName string, apiKey string) error {
	logger := log.WithFields(log.Fields{
		"func":       "ValidateCredential",
		"platform":   platformName,
		"request_id": ctx.Value("request_id"),
	})

	logger.Debug("Validating platform credential started")

	switch platformName {
	case "render":
		logger.Debug("Validating Render credentials")
		client := platform.NewRenderProvider(apiKey)
		err := client.VerifyCredentials(ctx)
		if err != nil {
			logger.WithError(err).Warn("Render credential validation failed")
		} else {
			logger.Debug("Render credential validated successfully")
		}
		return err
	case "vercel":
		logger.Warn("Vercel validation not implemented")
		// TODO: Implement Vercel validation
		return fmt.Errorf("vercel validation not implemented")
	default:
		logger.Warn("Unsupported platform specified")
		return fmt.Errorf("unsupported platform: %s", platformName)
	}
}
