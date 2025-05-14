package handler

import (
	"checkmate/api/internal/auth"
	"checkmate/api/internal/model"
	"checkmate/api/internal/service"
	"checkmate/api/internal/utils"
	"encoding/json"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func GetCredentials(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"handler":    "GetCredentials",
		"request_id": r.Context().Value("request_id"),
	})

	logger.Info("Getting credentials started")

	// get user ID from context
	userID, err := auth.GetUserFromRequest(r)
	if err != nil || userID == "" {
		logger.WithError(err).Warn("Unauthorized access attempt")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	logger = logger.WithField("user_id", userID)
	logger.Debug("User authenticated successfully")

	// get all credentials for user
	credentials, err := service.GetPlatformCredentials(r.Context(), userID)
	if err != nil {
		logger.WithError(err).Error("Failed to get platform credentials")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.WithField("credentials_count", len(credentials)).Debug("Retrieved credentials")

	// convert to safe credentials (without API keys)
	safeCredentials := utils.ConvertToSafeCredentials(credentials)

	// prepare and send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"credentials": safeCredentials,
	})

	logger.Info("Credentials successfully returned")
}

func CreateCredentials(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"handler":    "CreateCredentials",
		"request_id": r.Context().Value("request_id"),
	})

	logger.Info("Creating credential started")

	// get user ID from context
	userID, err := auth.GetUserFromRequest(r)
	if err != nil || userID == "" {
		logger.WithError(err).Warn("Unauthorized access attempt")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	logger = logger.WithField("user_id", userID)
	logger.Debug("User authenticated successfully")

	//parse request body
	var input model.PlatformCredentialInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.WithError(err).Warn("Failed to parse request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	logger = logger.WithField("platform", input.Platform)
	logger.Debug("Request body parsed successfully")

	// validate credential with platform before saving
	if err := service.ValidateCredential(r.Context(), input.Platform, input.APIKey); err != nil {
		logger.WithError(err).Warn("Credential validation failed")
		http.Error(w, "Invalid credential: "+err.Error(), http.StatusBadRequest)
		return
	}

	logger.Debug("Credential validated successfully")

	//create credential
	cred, err := service.CreatePlatformCredential(r.Context(), userID, &input)
	if err != nil {
		logger.WithError(err).Error("Failed to create platform credential")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.WithField("credential_id", cred.ID).Debug("Credential created successfully")

	// convert to safe credential (without API key)
	safeCred := utils.ConvertToSafeCredential(cred)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(safeCred)

	logger.WithField("credential_id", cred.ID).Info("New credential successfully created and returned")
}

func UpdateCredential(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"handler":    "UpdateCredential",
		"request_id": r.Context().Value("request_id"),
	})

	logger.Info("Updating credential started")

	// get user ID from context
	userID, err := auth.GetUserFromRequest(r)
	if err != nil || userID == "" {
		logger.WithError(err).Warn("Unauthorized access attempt")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	logger = logger.WithField("user_id", userID)
	logger.Debug("User authenticated successfully")

	// get credential ID from query params
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		logger.Warn("Missing credential ID in request")
		http.Error(w, "Missing credential ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.WithError(err).Warn("Invalid credential ID format")
		http.Error(w, "Invalid credential ID", http.StatusBadRequest)
		return
	}

	logger = logger.WithField("credential_id", id)
	logger.Debug("Credential ID parsed successfully")

	// parse request body
	var input model.PlatformCredentialInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.WithError(err).Warn("Failed to parse request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	logger = logger.WithField("platform", input.Platform)
	logger.Debug("Request body parsed successfully")

	// validate credential before updating
	if err := service.ValidateCredential(r.Context(), input.Platform, input.APIKey); err != nil {
		logger.WithError(err).Warn("Credential validation failed")
		http.Error(w, "Invalid credential: "+err.Error(), http.StatusBadRequest)
		return
	}

	logger.Debug("Credential validated successfully")

	// update credential
	if err := service.UpdatePlatformCredential(r.Context(), id, userID, &input); err != nil {
		logger.WithError(err).Error("Failed to update platform credential")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info("Credential successfully updated")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Credential updated successfully"}`))
}

func DeleteCredential(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"handler":    "DeleteCredential",
		"request_id": r.Context().Value("request_id"),
	})

	logger.Info("Deleting credential started")

	// get user ID from context
	userID, err := auth.GetUserFromRequest(r)
	if err != nil || userID == "" {
		logger.WithError(err).Warn("Unauthorized access attempt")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	logger = logger.WithField("user_id", userID)
	logger.Debug("User authenticated successfully")

	// get credential ID from query params
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		logger.Warn("Missing credential ID in request")
		http.Error(w, "Missing credential ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.WithError(err).Warn("Invalid credential ID format")
		http.Error(w, "Invalid credential ID", http.StatusBadRequest)
		return
	}

	logger = logger.WithField("credential_id", id)
	logger.Debug("Credential ID parsed successfully")

	// delete credential
	if err := service.DeletePlatformCredential(r.Context(), id, userID); err != nil {
		logger.WithError(err).Error("Failed to delete platform credential")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info("Credential successfully deleted")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Credential deleted successfully"}`))
}
