package handler

import (
	"checkmate/api/internal/auth"
	"checkmate/api/internal/service"
	"encoding/json"
	"net/http"
)

func GetDeployments(w http.ResponseWriter, r *http.Request) {
	//get the user id from context
	userID, err := auth.GetUserFromRequest(r)
	if err != nil || userID == ""{
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
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
