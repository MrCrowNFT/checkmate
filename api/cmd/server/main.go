package main

import (
	"checkmate/api/internal/auth"
	"checkmate/api/internal/handler"
	"checkmate/api/internal/storage"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current directory: %v", err)
	}
	log.Printf("Running from directory: %s", currentDir)
	storage.InitDb()

	err = godotenv.Load("./.env")
	if err != nil {
		//not necesary fail, since .env may not exists on production
		log.Printf("Warning: .env file not found: %v", err)
	}

	firebaseCredPath := os.Getenv("FIREBASE_CREDENTIALS_PATH")
	if firebaseCredPath == "" {
		firebaseCredPath = "../../internal/config/firebase-credentials.json" // default
		log.Printf("Using default Firebase credentials path: %s", firebaseCredPath)
	}

	err = auth.InitFirebase(firebaseCredPath)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}

	mux := http.NewServeMux()

	//endpoints
	mux.HandleFunc("/auth", auth.Authenticate(handler.GetCurrentUser))

	// Apply CORS middleware
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:1420"}, // Tauri default dev port
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := corsMiddleware.Handler(mux)

	server := &http.Server{
		Addr:    ":8080", //eventually add this to .env
		Handler: handler,
	}

	//start server in goroutine
	go func() {
		log.Printf("Server starting on port %s", server.Addr)
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	//GRACEFULL SHUTDOWN

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped gracefully")
}
