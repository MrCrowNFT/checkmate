package main

import (
	"checkmate/api/internal/auth"
	"checkmate/api/internal/handler"
	"checkmate/api/internal/storage"
	"checkmate/api/internal/utils"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

func main() {
	// init logger
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	if os.Getenv("ENV") == "production" {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(log.DebugLevel)
	}

	logger := log.WithFields(log.Fields{
		"func": "main",
	})

	// create package-level context key for request ID
	utils.InitRequestIDKey()

	// Debugging the working directory
	currentDir, err := os.Getwd()
	if err != nil {
		logger.WithError(err).Fatal("Failed to get current directory")
	}
	logger.WithField("directory", currentDir).Debug("Running from directory")

	// init database
	storage.InitDb()
	logger.Debug("Database initialized successfully")

	// get environment variables
	err = godotenv.Load("./.env")
	if err != nil {
		// Not necessary fail, since .env may not exist in production
		logger.WithError(err).Warn("Warning: .env file not found")
	} else {
		logger.Debug("Environment variables loaded successfully")
	}

	// get Firebase credentials path to init auth
	firebaseCredPath := os.Getenv("FIREBASE_CREDENTIALS_PATH")
	if firebaseCredPath == "" {
		firebaseCredPath = "../../internal/config/firebase-credentials.json" // default
		logger.WithField("path", firebaseCredPath).Debug("Using default Firebase credentials path")
	}

	err = auth.InitFirebase(firebaseCredPath)
	if err != nil {
		logger.WithError(err).Fatal("Failed to initialize Firebase")
	}
	logger.Debug("Firebase authentication initialized successfully")

	// Initialize encryption for platform credentials
	if err := utils.InitEncryption(); err != nil {
		logger.WithError(err).Fatal("Failed to initialize encryption")
	}
	logger.Debug("Encryption initialized successfully")

	mux := http.NewServeMux()

	// endpoints
	mux.HandleFunc("/", auth.AuthenticateWithRequestID(handler.GetCurrentUser))

	mux.HandleFunc("/deployments", auth.AuthenticateWithRequestID(handler.GetDeployments))
	mux.HandleFunc("/credentials", auth.AuthenticateWithRequestID(handler.GetCredentials))
	mux.HandleFunc("/credentials/new", auth.AuthenticateWithRequestID(handler.CreateCredentials))
	mux.HandleFunc("/credentials/update/:id", auth.AuthenticateWithRequestID(handler.UpdateCredential))
	mux.HandleFunc("/credentials/delete/:id", auth.AuthenticateWithRequestID(handler.DeleteCredential))

	logger.Debug("Routes registered successfully")

	// Apply CORS middleware
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:1420", "http://localhost:5173"}, // Tauri default dev port + current frontend
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := corsMiddleware.Handler(mux)
	logger.Debug("CORS middleware applied")

	server := &http.Server{
		Addr:    ":8080", // eventually add this to .env
		Handler: handler,
	}

	// Start server in goroutine
	go func() {
		logger.WithField("port", server.Addr).Info("Server starting")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Fatal("Server failed to start")
		}
	}()

	// GRACEFUL SHUTDOWN
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// wait for shutdown signal
	sig := <-quit
	logger.WithField("signal", sig.String()).Info("Shutdown signal received")

	// deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// shutdown server
	logger.Info("Shutting down server...")
	if err := server.Shutdown(ctx); err != nil {
		logger.WithError(err).Fatal("Server forced to shutdown")
	}

	logger.Info("Server stopped gracefully")
}
