package auth

import (
	"checkmate/api/internal/utils"
	"context"
	"fmt"
	"net/http"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

var (
	firebaseApp *firebase.App
	authClient  *auth.Client
)

// InitFirebase initializes the Firebase application and authentication client
func InitFirebase(credentialsPath string) error {
	logger := log.WithFields(log.Fields{
		"func": "InitFirebase",
		"path": credentialsPath,
	})

	logger.Debug("Initializing Firebase")

	opt := option.WithCredentialsFile(credentialsPath)

	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		logger.WithError(err).Error("Error initializing Firebase app")
		return fmt.Errorf("error initializing firebase app: %v", err)
	}

	firebaseApp = app
	logger.Debug("Firebase app initialized successfully")

	client, err := app.Auth(context.Background())
	if err != nil {
		logger.WithError(err).Error("Error getting auth client")
		return fmt.Errorf("error getting auth client: %v", err)
	}

	authClient = client
	logger.Debug("Firebase auth client initialized successfully")
	return nil
}

// Authenticate middleware to verify user authentication
func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := log.WithFields(log.Fields{
			"func":   "Authenticate",
			"path":   r.URL.Path,
			"method": r.Method,
		})

		logger.Debug("Processing authentication request")

		// Extract auth header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			logger.Warn("Authorization header missing")
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Check token format -> should have bearer (jwt)
		idToken := strings.TrimPrefix(authHeader, "Bearer ")
		if idToken == authHeader {
			logger.Warn("Invalid Authorization header format")
			http.Error(w, "Authorization header format must be Bearer {token}", http.StatusUnauthorized)
			return
		}

		// Verify token id
		token, err := authClient.VerifyIDToken(r.Context(), idToken)
		if err != nil {
			logger.WithError(err).Warn("Invalid token")
			http.Error(w, fmt.Sprintf("Invalid token: %v", err), http.StatusUnauthorized)
			return
		}

		logger.WithField("uid", token.UID).Debug("Authentication successful")

		ctx := context.WithValue(r.Context(), "user", token)
		ctx = context.WithValue(ctx, "uid", token.UID)

		// Call next handler with the enhanced context
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// AuthenticateWithRequestID adds request ID tracking to the authentication process
func AuthenticateWithRequestID(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Generate request ID
		requestID := utils.GenerateRequestID()

		// Add request ID to response headers for client-side tracking
		w.Header().Set("X-Request-ID", requestID)

		// Create context with request ID
		ctx := context.WithValue(r.Context(), utils.RequestIDKey, requestID)

		logger := log.WithFields(log.Fields{
			"func":       "AuthenticateWithRequestID",
			"path":       r.URL.Path,
			"method":     r.Method,
			"request_id": requestID,
		})

		logger.Debug("Processing authentication request with request ID")

		// Extract auth header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			logger.Warn("Authorization header missing")
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Check token format -> should have bearer (jwt)
		idToken := strings.TrimPrefix(authHeader, "Bearer ")
		if idToken == authHeader {
			logger.Warn("Invalid Authorization header format")
			http.Error(w, "Authorization header format must be Bearer {token}", http.StatusUnauthorized)
			return
		}

		// Verify token id
		token, err := authClient.VerifyIDToken(ctx, idToken)
		if err != nil {
			logger.WithError(err).Warn("Invalid token")
			http.Error(w, fmt.Sprintf("Invalid token: %v", err), http.StatusUnauthorized)
			return
		}

		logger.WithField("uid", token.UID).Debug("Authentication successful")

		ctx = context.WithValue(ctx, "user", token)
		ctx = context.WithValue(ctx, "uid", token.UID)

		// Call next handler with the enhanced context
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// GetUserFromRequest extracts the user ID from the request context
func GetUserFromRequest(r *http.Request) (string, error) {
	uid, ok := r.Context().Value("uid").(string)
	if !ok {
		return "", fmt.Errorf("user not found in request context")
	}
	return uid, nil
}

// GetTokenFromRequest extracts the token from the request context
func GetTokenFromRequest(r *http.Request) (*auth.Token, error) {
	token, ok := r.Context().Value("user").(*auth.Token)
	if !ok {
		return nil, fmt.Errorf("token not found in request context")
	}
	return token, nil
}

// GetRequestIDFromRequest extracts the request ID from the request context
func GetRequestIDFromRequest(r *http.Request) string {
	return utils.GetRequestIDFromContext(r.Context())
}
