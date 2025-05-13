package handler

import (
	"checkmate/api/internal/model"
	"checkmate/api/internal/service"
	"encoding/json"
	"net/http"
)

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
