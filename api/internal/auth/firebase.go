package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

var (
	firebaseApp *firebase.App
	authClient  *auth.Client
)

// initialize firebase
func InitFirebase(creadentialsPath string) error {
	opt := option.WithCredentialsFile(creadentialsPath)

	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return fmt.Errorf("Error initializing firebase app: %v", err)
	}

	firebaseApp = app

	client, err := app.Auth(context.Background())
	if err != nil {
		return fmt.Errorf("Error getting auth client: %v", err)
	}
	authClient = client
	return nil
}

func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//extract auth header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		//check token format -> should have bearer (jwt)
		idToken := strings.TrimPrefix(authHeader, "Bearer ")
		if idToken == authHeader {
			http.Error(w, "Authorization header format must be Bearer {token}", http.StatusUnauthorized)
			return
		}

		//verify token id
		token, err := authClient.VerifyIDToken(r.Context(), idToken)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid token: %v", err), http.StatusUnauthorized) // Added error details for debugging
			return
		}

		ctx := context.WithValue(r.Context(), "user", token)
		ctx = context.WithValue(ctx, "uid", token.UID)

		// call next handler with the enhanced context
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// extracts the user ID from the request context
func GetUserFromRequest(r *http.Request) (string, error) {
	uid, ok := r.Context().Value("uid").(string)
	if !ok {
		return "", fmt.Errorf("user not found in request context")
	}
	return uid, nil
}

// extracts the token from the request context
func GetTokenFromRequest(r *http.Request) (*auth.Token, error) {
	token, ok := r.Context().Value("user").(*auth.Token)
	if !ok {
		return nil, fmt.Errorf("token not found in request context")
	}
	return token, nil
}
