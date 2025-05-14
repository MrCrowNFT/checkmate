package handler

import (
	"checkmate/api/internal/auth"
	"checkmate/api/internal/service"
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func GetDeployments(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"handler":    "GetDeployments",
		"request_id": r.Context().Value("request_id"),
	})

	logger.Info("Getting deployments started")

	//get the user id from context
	userID, err := auth.GetUserFromRequest(r)
	if err != nil || userID == "" {
		logger.WithError(err).Warn("Unauthorized access attempt")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	logger = logger.WithField("user_id", userID)
	logger.Debug("User authenticated successfully")

	//get deployments
	deployments, err := service.GetAllUserDeployments(r.Context(), userID)
	if err != nil {
		logger.WithError(err).Error("Failed to retrieve user deployments")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.WithField("deployments_count", len(deployments)).Debug("Retrieved deployments")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	//* Note-> Remeber to decode on the frontend
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"deployments": deployments,
	}); err != nil {
		logger.WithError(err).Error("Failed to encode deployments response")
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	logger.Info("Deployments successfully returned")
}
