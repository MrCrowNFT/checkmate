package service

import (
	"checkmate/api/internal/model"
	"checkmate/api/internal/storage"
	"checkmate/api/internal/platform"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
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
	tx, err := storage.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() //roll back if not committed

	_, err = tx.ExecContext(ctx,
		"DELETE FROM deployment_cache WHERE platform_credential_id = ?",
		credentialID)
	if err != nil {
		return fmt.Errorf("failed to clear existing cache: %w", err)
	}

	//insert new deployments into cache
	now := time.Now()
	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO deployment_cache (
			id, platform_credential_id, name, status, url, 
			last_deployed_at, branch, service_type, framework, 
			last_updated_at, metadata
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()
	for _, dep := range deployments {
		// Convert metadata to JSON string
		metadataJSON, err := json.Marshal(dep.Metadata)
		if err != nil {
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
			return fmt.Errorf("failed to cache deployment %s: %w", dep.ID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

//get cached deployments for a platform
func GetCachedDeployments(ctx context.Context, credentialID int) ([]model.Deployment, time.Time, error) {
	query := `
		SELECT 
			id, name, status, url, last_deployed_at, branch, 
			service_type, framework, last_updated_at, metadata
		FROM deployment_cache
		WHERE platform_credential_id = ?
	`

	rows, err := storage.DB.QueryContext(ctx, query, credentialID)
	if err != nil {
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
				return nil, time.Time{}, fmt.Errorf("failed to unmarshal metadata: %w", err)
			}
			dep.Metadata = metadata
		} else {
			dep.Metadata = make(map[string]interface{})
		}

		dep.LastUpdatedAt = lastUpdatedAt
		deployments = append(deployments, dep)
	}

	return deployments, lastUpdatedAt, nil
}

// checks if cache exists for a credential
func CacheExists(ctx context.Context, credentialID int) (bool, time.Time, error) {
	query := `
		SELECT COUNT(*), MAX(last_updated_at) 
		FROM deployment_cache 
		WHERE platform_credential_id = ?
	`

	var count int
	var lastUpdated sql.NullTime

	err := storage.DB.QueryRowContext(ctx, query, credentialID).Scan(&count, &lastUpdated)
	if err != nil {
		return false, time.Time{}, fmt.Errorf("failed to check cache existence: %w", err)
	}

	// if no rows or null last_updated_at, cache doesn't exist
	if count == 0 || !lastUpdated.Valid {
		return false, time.Time{}, nil
	}

	return true, lastUpdated.Time, nil
}

// gets fresh deployments from cache or updates cache if stale
func GetFreshOrUpdateCache(ctx context.Context, cred *model.PlatformCredential) ([]model.Deployment, error) {
	// check if cache exists and is fresh
	exists, lastUpdated, err := CacheExists(ctx, cred.ID)
	if err != nil {
		return nil, err
	}

	if exists && IsCacheFresh(lastUpdated) {
		// cache is fresh, return cached data
		deployments, _, err := GetCachedDeployments(ctx, cred.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve cached deployments: %w", err)
		}
		return deployments, nil
	}

	// if cache doesn't exist or is stale, fetch fresh data from platform
	deployments, err := fetchDeploymentsFromPlatform(ctx, cred)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch deployments from platform: %w", err)
	}

	// update cache with fresh data
	err = StoreCachedDeployment(ctx, cred.ID, deployments)
	if err != nil {
		return nil, fmt.Errorf("failed to update cache: %w", err)
	}

	return deployments, nil
}

// fetches fresh deployment data from the platform
func fetchDeploymentsFromPlatform(ctx context.Context, cred *model.PlatformCredential) ([]model.Deployment, error) {
	switch cred.Platform {
	case "render":
		client := platform.NewRenderClient(cred.APIKey)
		deployments, err := client.GetServices(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch Render deployments: %w", err)
		}

		// Set PlatformCredential for each deployment
		for i := range deployments {
			// create copy to avoid modifying the original
			credCopy := *cred
			deployments[i].PlatformCredential = &credCopy
		}
		
		return deployments, nil
		
	case "vercel":
		// TODO: Implement Vercel client
		return nil, fmt.Errorf("vercel platform not implemented")
		
	default:
		return nil, fmt.Errorf("unsupported platform: %s", cred.Platform)
	}
}

// gets deployments for all user credentials this is the main function here
//workflow-> first get all platform credentials associated to the user id -> check if cache is fresh or stale
//-> if is fresh it returns the cache
//-> if not, fetch the data from the paltform ->case render/case vercel/etc-> update the cache and return it
//-> append each deployment to the array and return it
func GetAllUserDeployments(ctx context.Context, userID string) ([]model.Deployment, error) {
	// get all credentials for the user
	creds, err := GetPlatformCredentials(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user credentials: %w", err)
	}

	// Collect all deployments from all credentials
	var allDeployments []model.Deployment
	
	for _, cred := range creds {
		deployments, err := GetFreshOrUpdateCache(ctx, &cred)
		if err != nil {
			// Log error but continue with other credentials
			fmt.Printf("Error fetching deployments for credential %d: %v\n", cred.ID, err)
			continue
		}
		
		allDeployments = append(allDeployments, deployments...)
	}

	return allDeployments, nil
}
