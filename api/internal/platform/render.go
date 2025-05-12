package platform

import (
	"checkmate/api/internal/model"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
//render returns an array of this
type RenderService struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	Type              string    `json:"type"`
	ServiceDetails    *ServiceDetails `json:"serviceDetails"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
	Suspended         string    `json:"suspended"`
	URL               string    `json:"url"`
	Branch            string    `json:"branch"`
	Status            string    `json:"status"`
	AutoDeploy        bool      `json:"autoDeploy"`
	LastSuccessfulDeployAt *time.Time `json:"lastSuccessfulDeployAt"`
}

type ServiceDetails struct {
	BuildCommand         string `json:"buildCommand"`
	StartCommand         string `json:"startCommand"`
	Env                  string `json:"env"`
	PlanID               string `json:"planId"`
	Region               string `json:"region"`
	Pullrequests         string `json:"pullrequests"`
	NumInstances         int    `json:"numInstances"`
	Domains              []string `json:"domains"`
	DetectedFramework    string `json:"detectedFramework"`
}