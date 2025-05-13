package handler

import (
	"checkmate/api/internal/model"
	"checkmate/api/internal/service"
	"encoding/json"
	"net/http"
	"strconv"
)
//todo this should return the credentials but without the apikey for common sense reasons
func GetCredentials(w http.ResponseWriter, r *http.Request) {
	// get user ID from context
	userID := r.Context().Value("userId").(string)
	
	// get all credentials for user
	credentials, err := service.GetPlatformCredentials(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// prepare and send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"credentials": credentials,
	})
}


func CreateCredentials(w http.ResponseWriter, r *http.Request) {
	//todo check if this user id value is right
	userID := r.Context().Value("userId").(string)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	//parse request body
	var input model.PlatformCredentialInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// validate credential with platform before saving
	if err := service.ValidateCredential(r.Context(), input.Platform, input.APIKey); err != nil {
		http.Error(w, "Invalid credential: "+err.Error(), http.StatusBadRequest)
		return
	}

	//create credential
	cred, err := service.CreatePlatformCredential(r.Context(), userID, &input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cred)//*note the cred has the api key, is it okay to send it back like this?
}


func UpdateCredential(w http.ResponseWriter, r *http.Request) {
	// get user ID from context
	userID := r.Context().Value("userId").(string)
	
	// get credential ID from query params
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing credential ID", http.StatusBadRequest)
		return
	}
	
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid credential ID", http.StatusBadRequest)
		return
	}
	
	// Parse request body
	var input model.PlatformCredentialInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Validate credential with platform before updating
	if err := service.ValidateCredential(r.Context(), input.Platform, input.APIKey); err != nil {
		http.Error(w, "Invalid credential: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	// Update credential
	if err := service.UpdatePlatformCredential(r.Context(), id, userID, &input); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Credential updated successfully"}`))
}
