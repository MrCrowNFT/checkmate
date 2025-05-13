package handler

import (
	"checkmate/api/internal/service"
	"encoding/json"
	"net/http"
)

func GetDeployments(w http.ResponseWriter, r *http.Request) {
	//get the user id from context
	//todo check if this user id value is right
	userID := r.Context().Value("userId").(string)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	//get deployments
	deployments, err := service.GetAllUserDeployments(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	//* Note-> Remeber to decode on the frontend
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"deployments": deployments,
	}); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
