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

type RenderClient struct {
	apiKey string
	client *http.Client
}

func NewRenderClient(apiKey string) *RenderClient {
	return &RenderClient{
		apiKey: apiKey,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// render returns an array of this
type RenderService struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Type         string    `json:"type"`
	Branch       string    `json:"branch"`
	Suspended    string    `json:"suspended"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	AutoDeploy   string    `json:"autoDeploy"`
	Repo         string    `json:"repo"`
	DashboardURL string    `json:"dashboardUrl"`

	ServiceDetails struct {
		BuildCommand string `json:"buildCommand"`
		PublishPath  string `json:"publishPath"`
		URL          string `json:"url"`
		BuildPlan    string `json:"buildPlan"`
		ParentServer *struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"parentServer,omitempty"`
	} `json:"serviceDetails"`
}

type RenderServiceResponse struct {
	Service RenderService `json:"service"`
	Cursor  string        `json:"cursor,omitempty"`
}

type RenderServicesResponse struct {
	Services []RenderServiceResponse `json:"services"`
	Cursor   string                  `json:"cursor,omitempty"`
}

// verify valid api key
func (c *RenderClient) VerifyCredentials(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", renderAPIBaseURL+"services", nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
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

func (c *RenderClient) GetServices(ctx context.Context) ([]model.Deployment, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", renderAPIBaseURL+"services", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
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
	var services []RenderService
	if err := json.NewDecoder(resp.Body).Decode(&services); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	deployments := make([]model.Deployment, 0, len(services))
	for _, service := range services {
		//determine status
		status := determineDeploymentStatus(service)

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

		//append as deployment
		deployments = append(deployments, model.Deployment{
			ID:             service.ID,
			Name:           service.Name,
			Status:         status,
			URL:            service.ServiceDetails.URL,
			LastDeployedAt: lastDeployed,
			Branch:         service.Branch,
			ServiceType:    service.Type,
			Framework:      inferFrameworkFromRepo(service.Repo),
			LastUpdatedAt:  service.UpdatedAt,
			Metadata:       metadata,
		})
	}
	return deployments, nil

}

// todo there are more status in render, need to check them out
func determineDeploymentStatus(service RenderService) model.DeploymentStatus {
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

func inferFrameworkFromRepo(repoURL string) string {
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
