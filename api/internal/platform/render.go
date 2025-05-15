package platform

import (
	"checkmate/api/internal/model"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const renderAPIBaseURL = "https://api.render.com/v1"

// implements operations for the Render platform
type RenderProvider struct {
	client *model.RenderClient
}

func NewRenderProvider(apiKey string) *RenderProvider {
	return &RenderProvider{
		client: &model.RenderClient{
			ApiKey: apiKey,
			Client: &http.Client{
				Timeout: 30 * time.Second,
			},
		},
	}
}

// verify valid api key
func (p *RenderProvider) VerifyCredentials(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", renderAPIBaseURL+"/services", nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+p.client.ApiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := p.client.Client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("invalid API key")
	} else if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("received non-OK response: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (p *RenderProvider) GetServices(ctx context.Context) ([]model.Deployment, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", renderAPIBaseURL+"/services", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+p.client.ApiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := p.client.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("received non-OK response: %d, body: %s", resp.StatusCode, string(body))
	}

	//decode response into renderServiceResponses array so that we can loop
	//and append each one as a mode.deployment into new array
	var serviceResponses []model.RenderServiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&serviceResponses); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	deployments := make([]model.Deployment, 0, len(serviceResponses))
	for _, response := range serviceResponses {
		service := response.Service
		//determine status
		status := p.determineDeploymentStatus(service)

		metadata := map[string]interface{}{
			"type":         service.Type,
			"autoDeploy":   service.AutoDeploy,
			"suspended":    service.Suspended,
			"createdAt":    service.CreatedAt,
			"updatedAt":    service.UpdatedAt,
			"repo":         service.Repo,
			"dashboardUrl": service.DashboardURL,
			"buildPlan":    service.ServiceDetails.BuildPlan,
		}

		//add parent server if exists
		if service.ServiceDetails.ParentServer != nil {
			metadata["parentServerId"] = service.ServiceDetails.ParentServer.ID
			metadata["parentServerName"] = service.ServiceDetails.ParentServer.Name
		}

		// calculate last deployed time (approximation based on updatedAt) for caching purpuses
		var lastDeployed *time.Time
		if service.Status == "live" || service.UpdatedAt != service.CreatedAt {
			lastDeployed = &service.UpdatedAt
		}

		//todo missing the PlatformCredentialID
		//append as deployment
		deployments = append(deployments, model.Deployment{
			ID:             service.ID,
			Name:           service.Name,
			Status:         status,
			URL:            service.ServiceDetails.URL,
			LastDeployedAt: lastDeployed,
			Branch:         service.Branch,
			ServiceType:    service.Type,
			Framework:      p.inferFrameworkFromRepo(service.Repo),
			LastUpdatedAt:  service.UpdatedAt,
			Metadata:       metadata,
		})
	}
	return deployments, nil
}

// todo there are more status in render, need to check them out
func (p *RenderProvider) determineDeploymentStatus(service model.RenderService) model.DeploymentStatus {
	switch strings.ToLower(service.Status) {
	case "live", "up":
		return model.DeploymentStatusLive
	case "suspended":
		return model.DeploymentStatusCanceled
	case "deploying", "build":
		return model.DeploymentStatusDeploying
	case "failed", "error":
		return model.DeploymentStatusFailed
	default:
		return model.DeploymentStatusUnknown
	}
}

func (p *RenderProvider) inferFrameworkFromRepo(repoURL string) string {
	if repoURL == "" {
		return ""
	}

	repoURL = strings.ToLower(repoURL)
	frameworks := map[string]string{
		"react":   "react",
		"vue":     "vue",
		"angular": "angular",
		"nextjs":  "next.js",
		"nuxtjs":  "nuxt.js",
		"gatsby":  "gatsby",
		"svelte":  "svelte",
		"remix":   "remix",
		"astro":   "astro",
		"express": "express",
		"nestjs":  "nest.js",
		"flask":   "flask",
		"django":  "django",
		"fastapi": "fastapi",
	}

	for framework, name := range frameworks {
		if strings.Contains(repoURL, framework) {
			return name
		}
	}
	return ""
}
