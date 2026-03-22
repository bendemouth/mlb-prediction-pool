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

	// Initialize S3 client
	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		awsRegion = "us-east-1"
	}

	s3Client, err := handlers.NewS3Client(ctx, awsRegion)
	if err != nil {
		log.Fatal("Failed to initialize S3 client:", err)
	}

	// Create handlers
	h := handlers.NewHandler(db, *s3Client)

	// Create public server and routes
	publicMux := http.NewServeMux()

	publicMux.HandleFunc("/health", h.HandleHealthCheck)
	publicMux.HandleFunc("/leaderboard", h.GetLeaderboard)

	// Create protected server and routes
	protectedMux := http.NewServeMux()

	// Predictions endpoints
	protectedMux.HandleFunc("/predictions", h.GetPredictionsByUser)
	protectedMux.HandleFunc("/predictions/create", h.CreatePrediction)
	protectedMux.HandleFunc("/predictions/batchCreate", h.CreateBulkPredictions)
	protectedMux.HandleFunc("/predictions/game", h.GetPredictionsByGame)

	// User endpoints
	protectedMux.HandleFunc("/users/create", h.HandleCreateUser)
	protectedMux.HandleFunc("/users", h.HandleGetUser)
	protectedMux.HandleFunc("/users/listUsers", h.HandleListUsers)
	protectedMux.HandleFunc("/users/stats", h.HandleGetUserStats)

	// Games endpoints
	protectedMux.HandleFunc("/games/upcoming", h.GetUpcomingGamesSummary)

	// Model endpoints
	protectedMux.HandleFunc("/models/submitModel", h.UploadModelHandler)
	protectedMux.HandleFunc("/models", h.GetUserModelsHandler)
	protectedMux.HandleFunc("/models/delete/", h.DeleteModelHandler)
	protectedMux.HandleFunc("/models/", h.GetModelHandler)

	protectedHandler := middleware.Auth(protectedMux)

	mainMux := http.NewServeMux()
	mainMux.Handle("/health", publicMux)
	mainMux.Handle("/leaderboard", publicMux)
	mainMux.Handle("/", protectedHandler)

	// Add middleware
	handler := middleware.Logger(
		middleware.CORS(
			middleware.Recovery(mainMux),
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
