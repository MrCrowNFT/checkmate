package model

import (
	"time"
)

type DeploymentStatus string

const (
	DeploymentStatusLive      DeploymentStatus = "live"
	DeploymentStatusDeploying DeploymentStatus = "deploying"
	DeploymentStatusCanceled  DeploymentStatus = "canceled"
	DeploymentStatusFailed    DeploymentStatus = "failed"
	DeploymentStatusUnknown   DeploymentStatus = "unknown"
)

type Deployment struct {
	ID                   string                 `json:"id"`
	PlatformCredentialID int                    `json:"platformCredentialID"`
	Name                 string                 `json:"name"`
	Status               DeploymentStatus       `json:"status"`
	URL                  string                 `json:"url"`
	LastDeployedAt       *time.Time             `json:"lastDeployedAt"`
	Branch               string                 `json:"branch"`
	ServiceType          string                 `json:"serviceType"`
	Framework            string                 `json:"framework"`
	LastUpdatedAt        time.Time              `json:"lastUpdatedAt"`
	Metadata             map[string]interface{} `json:"metadata"`
}
