package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bendemouth/mlb-prediction-pool/internal/database"
	"github.com/bendemouth/mlb-prediction-pool/internal/handlers"
	"github.com/bendemouth/mlb-prediction-pool/internal/middleware"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	ctx := context.Background()

	// Initialize database with retries
	var db *database.DB
	var err error
	maxRetries := 10
	retryDelay := 2 * time.Second

	log.Println("Attempting to connect to database...")

	for i := 0; i < maxRetries; i++ {
		db, err = database.NewDBFromEnv(ctx)
		if err == nil {
			// Test the connection with health check
			healthErr := db.HealthCheck(ctx)
			if healthErr == nil {
				log.Println("Database connection established successfully")
				break
			}
			// Log the actual error
			log.Printf("Database health check failed (attempt %d/%d): %v",
				i+1, maxRetries, healthErr)
			err = healthErr
		} else {
			log.Printf("Database connection failed (attempt %d/%d): %v",
				i+1, maxRetries, err)
		}

		if i < maxRetries-1 {
			log.Printf("Retrying in %v...", retryDelay)
			time.Sleep(retryDelay)
		}
	}

	if err != nil {
		log.Fatal("Failed to connect to database after retries:", err)
	}

	defer db.Close()

	log.Println("Database connection established")

	// Create handlers
	h := handlers.NewHandler(db)

	// Create public server and routes
	mux := http.NewServeMux()

	mux.HandleFunc("/health", h.HandleHealthCheck)
	mux.HandleFunc("/leaderboard", h.GetLeaderboard)

	// Create protected server and routes
	authenticatedMux := http.NewServeMux()

	// Predictions endpoints
	authenticatedMux.HandleFunc("/predictions", h.GetPredictionsByUser)
	authenticatedMux.HandleFunc("/predictions/create", h.CreatePrediction)
	authenticatedMux.HandleFunc("/predictions/batchCreate", h.CreateBulkPredictions)
	authenticatedMux.HandleFunc("/predictions/game", h.GetPredictionsByGame)

	// User endpoints
	authenticatedMux.HandleFunc("/users/create", h.HandleCreateUser)
	authenticatedMux.HandleFunc("/users", h.HandleGetUser)
	authenticatedMux.HandleFunc("/users/listUsers", h.HandleListUsers)
	authenticatedMux.HandleFunc("/users/stats", h.HandleGetUserStats)

	// Add middleware
	handler := middleware.Logger(
		middleware.CORS(
			middleware.Recovery(mux),
		),
	)

	// Server config
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed:", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
