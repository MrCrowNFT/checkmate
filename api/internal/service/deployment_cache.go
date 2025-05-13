package service

import (
	"checkmate/api/internal/model"
	"checkmate/api/internal/storage"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

const (
	CacheTTL = 30 //seconds that the cached ployments are considered fresh
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
