package handler

import (
	"checkmate/api/internal/auth"
	"checkmate/api/internal/model"
	"checkmate/api/internal/service"
	"encoding/json"
	"log"
	"net/http"
)

// get current user profile
func GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	log.Println("GetCurrentUser handler called")

	// get user ID from request
	userID, err := auth.GetUserFromRequest(r)
	if err != nil {
		log.Printf("Error getting user from request: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	log.Printf("User ID from request: %s", userID)

	// get user from the database
	user, err := service.GetUserById(r.Context(), userID)
	if err != nil {
		log.Printf("Error retrieving user: %v", err)
		http.Error(w, "Error retrieving user", http.StatusInternalServerError)
		return
	}

	// if user doesn't exist, create a new record
	if user == nil {
		log.Println("User not found in database. Creating new user...")

		// get Firebase token to extract user info
		token, err := auth.GetTokenFromRequest(r)
		if err != nil {
			log.Printf("Error getting token: %v", err)
			http.Error(w, "Error getting token", http.StatusInternalServerError)
			return
		}

		log.Printf("Token retrieved for user %s", userID)

		// create new user model with Firebase user data
		user = &model.User{
			ID: userID,
		}

		// get email if available
		if email, ok := token.Claims["email"].(string); ok {
			user.Email = email
		} else {
			log.Println("Warning: No email found in token claims")
			// default email to avoid database constraints
			user.Email = userID + "@unknown.com"
		}

		// get display name if available
		if name, ok := token.Claims["name"].(string); ok {
			user.DisplayName = name
		} else {
			log.Println("No display name found in token claims")
		}

		// Save new user
		err = service.CreateUser(r.Context(), user)
		if err != nil {
			log.Printf("Error creating user: %v", err)
			http.Error(w, "Error creating user", http.StatusInternalServerError)
			return
		}

		log.Println("User successfully created")
	} else {
		log.Printf("User found in database: %s", user.Email)
	}

	// return the user data
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("Error encoding user to JSON: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	log.Println("User data successfully returned")
}
