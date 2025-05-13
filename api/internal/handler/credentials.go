package handler

import (
	"checkmate/api/internal/model"
	"checkmate/api/internal/service"
	"checkmate/api/internal/auth"
	"checkmate/api/internal/utils"
	"encoding/json"
	"net/http"
	"strconv"
)

func GetCredentials(w http.ResponseWriter, r *http.Request) {
	// get user ID from context
	userID, err := auth.GetUserFromRequest(r)
	if err != nil || userID == ""{
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	
	// get all credentials for user
	credentials, err := service.GetPlatformCredentials(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// convert to safe credentials (without API keys)
	safeCredentials := utils.ConvertToSafeCredentials(credentials)
	
	// prepare and send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"credentials": safeCredentials,
	})
}


func CreateCredentials(w http.ResponseWriter, r *http.Request) {
	// get user ID from context
	userID, err := auth.GetUserFromRequest(r)
	if err != nil || userID == ""{
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
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

	// convert to safe credential (without API key)
	safeCred := utils.ConvertToSafeCredential(cred)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(safeCred)//*note the cred has the api key, is it okay to send it back like this?
}


func UpdateCredential(w http.ResponseWriter, r *http.Request) {
	// get user ID from context
	userID, err := auth.GetUserFromRequest(r)
	if err != nil || userID == ""{
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	
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
	
	// parse request body
	var input model.PlatformCredentialInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// validate credential before updating
	if err := service.ValidateCredential(r.Context(), input.Platform, input.APIKey); err != nil {
		http.Error(w, "Invalid credential: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	// update credential
	if err := service.UpdatePlatformCredential(r.Context(), id, userID, &input); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Credential updated successfully"}`))
}

func DeleteCredential(w http.ResponseWriter, r *http.Request) {
	// get user ID from context
	userID, err := auth.GetUserFromRequest(r)
	if err != nil || userID == ""{
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

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
	
	// delete credential
	if err := service.DeletePlatformCredential(r.Context(), id, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Credential deleted successfully"}`))
}
