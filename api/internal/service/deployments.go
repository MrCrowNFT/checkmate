package service

import (
	"checkmate/api/internal/model"
	"checkmate/api/internal/platform"
	"checkmate/api/internal/storage"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	CacheTTL = 30 //seconds that the cached deployments are considered fresh
)

// checks if cached data is still valid (less than CacheTTL seconds old)
func IsCacheFresh(lastUpdatedAt time.Time) bool {
	return time.Since(lastUpdatedAt) < time.Duration(CacheTTL)*time.Second
}

// save deployment data on cache
func StoreCachedDeployment(ctx context.Context, credentialID int, deployments []model.Deployment) error {
	logger := log.WithFields(log.Fields{
		"func":              "StoreCachedDeployment",
		"credential_id":     credentialID,
		"deployments_count": len(deployments),
		"request_id":        ctx.Value("request_id"),
	})

	logger.Debug("Storing deployments in cache started")

	tx, err := storage.DB.BeginTx(ctx, nil)
	if err != nil {
		logger.WithError(err).Error("Failed to begin transaction")
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() //roll back if not committed

	_, err = tx.ExecContext(ctx,
		"DELETE FROM deployment_cache WHERE platform_credential_id = ?",
		credentialID)
	if err != nil {
		logger.WithError(err).Error("Failed to clear existing cache")
		return fmt.Errorf("failed to clear existing cache: %w", err)
	}

	logger.Debug("Cleared existing cache entries")

	//insert new deployments into cache
	now := time.Now() //we need this to insert in the last_updated_at

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO deployment_cache (
			id, platform_credential_id, name, status, url, 
			last_deployed_at, branch, service_type, framework, 
			last_updated_at, metadata
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		logger.WithError(err).Error("Failed to prepare statement")
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, dep := range deployments {
		// Convert metadata to JSON string
		metadataJSON, err := json.Marshal(dep.Metadata)
		if err != nil {
			logger.WithError(err).Error("Failed to marshal metadata")
			return fmt.Errorf("failed to marshal metadata: %w", err)
		}

		// handle null last_deployed_at
		var lastDeployedAt sql.NullTime
		if dep.LastDeployedAt != nil {
			lastDeployedAt = sql.NullTime{
				Time:  *dep.LastDeployedAt,
				Valid: true,
			}
		}

		_, err = stmt.ExecContext(ctx,
			dep.ID, credentialID, dep.Name, string(dep.Status), dep.URL,
			lastDeployedAt, dep.Branch, dep.ServiceType, dep.Framework,
			now, metadataJSON)
		if err != nil {
			logger.WithFields(log.Fields{
				"deployment_id":   dep.ID,
				"deployment_name": dep.Name,
			}).WithError(err).Error("Failed to cache deployment")
			return fmt.Errorf("failed to cache deployment %s: %w", dep.ID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		logger.WithError(err).Error("Failed to commit transaction")
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	logger.Info("Successfully stored deployments in cache")
	return nil
}

// get cached deployments for a platform
func GetCachedDeployments(ctx context.Context, credentialID int) ([]model.Deployment, time.Time, error) {
	logger := log.WithFields(log.Fields{
		"func":          "GetCachedDeployments",
		"credential_id": credentialID,
		"request_id":    ctx.Value("request_id"),
	})

	logger.Debug("Getting cached deployments started")

	query := `
		SELECT 
			id, name, status, url, last_deployed_at, branch, 
			service_type, framework, last_updated_at, metadata
		FROM deployment_cache
		WHERE platform_credential_id = ?
	`

	rows, err := storage.DB.QueryContext(ctx, query, credentialID)
	if err != nil {
		logger.WithError(err).Error("Failed to query cached deployments")
		return nil, time.Time{}, fmt.Errorf("failed to query cached deployments: %w", err)
	}
	defer rows.Close()

	var deployments []model.Deployment
	var lastUpdatedAt time.Time

	for rows.Next() {
		var dep model.Deployment
		var status string
		var metadataJSON string
		var lastDeployedAt sql.NullTime

		err := rows.Scan(
			&dep.ID, &dep.Name, &status, &dep.URL, &lastDeployedAt, &dep.Branch,
			&dep.ServiceType, &dep.Framework, &lastUpdatedAt, &metadataJSON,
		)
		if err != nil {
			logger.WithError(err).Error("Failed to scan deployment row")
			return nil, time.Time{}, fmt.Errorf("failed to scan deployment row: %w", err)
		}

		// Convert status string to DeploymentStatus
		dep.Status = model.DeploymentStatus(status)

		// Convert NullTime to *time.Time
		if lastDeployedAt.Valid {
			dep.LastDeployedAt = &lastDeployedAt.Time
		}

		// Parse metadata JSON
		if metadataJSON != "" {
			var metadata map[string]interface{}
			if err := json.Unmarshal([]byte(metadataJSON), &metadata); err != nil {
				logger.WithFields(log.Fields{
					"deployment_id":   dep.ID,
					"deployment_name": dep.Name,
				}).WithError(err).Error("Failed to unmarshal metadata")
				return nil, time.Time{}, fmt.Errorf("failed to unmarshal metadata: %w", err)
			}
			dep.Metadata = metadata
		} else {
			dep.Metadata = make(map[string]interface{})
		}

		dep.LastUpdatedAt = lastUpdatedAt
		deployments = append(deployments, dep)
	}

	logger.WithFields(log.Fields{
		"deployments_count": len(deployments),
		"last_updated_at":   lastUpdatedAt,
	}).Debug("Retrieved cached deployments successfully")
	return deployments, lastUpdatedAt, nil
}

// checks if cache exists for a credential
func CacheExists(ctx context.Context, credentialID int) (bool, time.Time, error) {
	logger := log.WithFields(log.Fields{
		"func":          "CacheExists",
		"credential_id": credentialID,
		"request_id":    ctx.Value("request_id"),
	})

	logger.Debug("Checking if cache exists")

	//we get the oldest last_updated_at deployment that has the platform id
	query := `
        SELECT COUNT(*), MAX(last_updated_at) 
        FROM deployment_cache 
        WHERE platform_credential_id = ?
    `

	var count int
	var lastUpdatedStr sql.NullString // Change to NullString instead of NullTime

	err := storage.DB.QueryRowContext(ctx, query, credentialID).Scan(&count, &lastUpdatedStr)
	if err != nil {
		logger.WithError(err).Error("Failed to check cache existence")
		return false, time.Time{}, fmt.Errorf("failed to check cache existence: %w", err)
	}

	// if no rows or null last_updated_at, cache doesn't exist
	if count == 0 || !lastUpdatedStr.Valid {
		logger.Debug("Cache does not exist")
		return false, time.Time{}, nil
	}

	// Parse the string into a time.Time
	lastUpdated, err := time.Parse(time.RFC3339, lastUpdatedStr.String)
	if err != nil {
		logger.WithError(err).Error("Failed to parse last_updated_at timestamp")
		return false, time.Time{}, fmt.Errorf("failed to parse timestamp: %w", err)
	}

	logger.WithFields(log.Fields{
		"cache_exists":    true,
		"last_updated_at": lastUpdated,
	}).Debug("Cache check complete")
	return true, lastUpdated, nil
}

// gets fresh deployments from cache or updates cache if stale
func GetFreshOrUpdateCache(ctx context.Context, cred *model.PlatformCredential) ([]model.Deployment, error) {
	logger := log.WithFields(log.Fields{
		"func":          "GetFreshOrUpdateCache",
		"credential_id": cred.ID,
		"platform":      cred.Platform,
		"request_id":    ctx.Value("request_id"),
	})

	logger.Debug("Getting fresh deployments or updating cache started")

	// check if cache exists and is fresh
	exists, lastUpdated, err := CacheExists(ctx, cred.ID)
	if err != nil {
		logger.WithError(err).Error("Failed to check if cache exists")
		return nil, err
	}

	if exists && IsCacheFresh(lastUpdated) {
		// cache is fresh, return cached data
		logger.WithField("last_updated_at", lastUpdated).Debug("Cache is fresh, using cached data")
		deployments, _, err := GetCachedDeployments(ctx, cred.ID)
		if err != nil {
			logger.WithError(err).Error("Failed to retrieve cached deployments")
			return nil, fmt.Errorf("failed to retrieve cached deployments: %w", err)
		}
		return deployments, nil
	}

	// if cache doesn't exist or is stale, fetch fresh data from platform
	logger.Debug("Cache is stale or doesn't exist, fetching fresh data")
	deployments, err := fetchDeploymentsFromPlatform(ctx, cred)
	if err != nil {
		logger.WithError(err).Error("Failed to fetch deployments from platform")
		return nil, fmt.Errorf("failed to fetch deployments from platform: %w", err)
	}

	// update cache with fresh data
	logger.Debug("Updating cache with fresh data")
	err = StoreCachedDeployment(ctx, cred.ID, deployments)
	if err != nil {
		logger.WithError(err).Error("Failed to update cache")
		return nil, fmt.Errorf("failed to update cache: %w", err)
	}

	logger.WithField("deployments_count", len(deployments)).Info("Successfully retrieved fresh deployments")
	return deployments, nil
}

// fetches fresh deployment data from the platform
func fetchDeploymentsFromPlatform(ctx context.Context, cred *model.PlatformCredential) ([]model.Deployment, error) {
	logger := log.WithFields(log.Fields{
		"func":          "fetchDeploymentsFromPlatform",
		"credential_id": cred.ID,
		"platform":      cred.Platform,
		"request_id":    ctx.Value("request_id"),
	})

	logger.Debug("Fetching deployments from platform started")

	switch cred.Platform {
	case "render":
		logger.Debug("Fetching deployments from Render")
		client := platform.NewRenderProvider(cred.APIKey)
		deployments, err := client.GetServices(ctx)
		if err != nil {
			logger.WithError(err).Error("Failed to fetch Render deployments")
			return nil, fmt.Errorf("failed to fetch Render deployments: %w", err)
		}

		// Set PlatformCredentialID for each deployment
		for i := range deployments {
			// create copy to avoid modifying the original
			deployments[i].PlatformCredentialID = cred.ID
		}

		logger.WithField("deployments_count", len(deployments)).Debug("Successfully fetched Render deployments")
		return deployments, nil

	case "vercel":
		logger.Warn("Vercel platform not implemented")
		// TODO: Implement Vercel client
		return nil, fmt.Errorf("vercel platform not implemented")

	default:
		logger.Warn("Unsupported platform specified")
		return nil, fmt.Errorf("unsupported platform: %s", cred.Platform)
	}
}

// gets deployments for all user credentials this is the main function here
// workflow-> first get all platform credentials associated to the user id -> check if cache is fresh or stale
// -> if is fresh it returns the cache
// -> if not, fetch the data from the paltform ->case render/case vercel/etc-> update the cache and return it
// -> append each deployment to the array and return it
func GetAllUserDeployments(ctx context.Context, userID string) ([]model.Deployment, error) {
	logger := log.WithFields(log.Fields{
		"func":       "GetAllUserDeployments",
		"user_id":    userID,
		"request_id": ctx.Value("request_id"),
	})

	logger.Info("Getting all user deployments started")

	// get all credentials for the user
	creds, err := GetPlatformCredentials(ctx, userID)
	if err != nil {
		logger.WithError(err).Error("Failed to get user credentials")
		return nil, fmt.Errorf("failed to get user credentials: %w", err)
	}

	logger.WithField("credentials_count", len(creds)).Debug("Retrieved user credentials")

	// Collect all deployments from all credentials
	var allDeployments []model.Deployment

	//all deployment assosiated to one credential should have
	//the same platform_credential_id
	for _, cred := range creds {
		credLogger := logger.WithFields(log.Fields{
			"credential_id": cred.ID,
			"platform":      cred.Platform,
		})

		credLogger.Debug("Processing credential")
		deployments, err := GetFreshOrUpdateCache(ctx, &cred)
		if err != nil {
			// Log error but continue with other credentials
			credLogger.WithError(err).Warn("Error fetching deployments for credential, continuing with others")
			fmt.Printf("Error fetching deployments for credential %d: %v\n", cred.ID, err)
			continue
		}

		credLogger.WithField("deployments_count", len(deployments)).Debug("Retrieved deployments for credential")
		allDeployments = append(allDeployments, deployments...)
	}

	logger.WithField("total_deployments", len(allDeployments)).Info("Successfully retrieved all user deployments")
	return allDeployments, nil
}
