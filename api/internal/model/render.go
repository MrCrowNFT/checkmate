package model

import (
	"net/http"
	"time"
)

type RenderClient struct {
	ApiKey string
	Client *http.Client
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
