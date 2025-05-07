package handler

import (
	"checkmate/api/internal/auth"
	"checkmate/api/internal/model"
	"checkmate/api/internal/service"
	"net/http"
	"encoding/json"
)

//get current user profile
func GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// get user ID from request 
	userID, err := auth.GetUserFromRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// get user from the database
	user, err := service.GetUserById(r.Context(), userID)
	if err != nil {
		http.Error(w, "Error retrieving user", http.StatusInternalServerError)
		return
	}

	// if user doesn't exist, create a new record
	if user == nil {
		// get Firebase token to extract user info
		token, err := auth.GetTokenFromRequest(r)
		if err != nil {
			http.Error(w, "Error getting token", http.StatusInternalServerError)
			return
		}

		// create new user model with Firebase user data
		user = &model.User{
			ID:    userID,
			Email: token.Claims["email"].(string),
		}

		// get display name if available
		if name, ok := token.Claims["name"].(string); ok {
			user.DisplayName = name
		}

		// Save new user
		err = service.CreateUser(r.Context(), user)
		if err != nil {
			http.Error(w, "Error creating user", http.StatusInternalServerError)
			return
		}
	}

	// return the user data
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)//encode user data
}


