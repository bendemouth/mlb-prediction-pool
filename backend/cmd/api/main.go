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

	// Initialize database
	db, err := database.NewDBFromEnv(ctx)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	log.Println("Database connection established")

	// Create handlers
	h := handlers.NewHandler(db)

	// Setup routes
	mux := http.NewServeMux()

	mux.HandleFunc("/health", h.HandleHealthCheck)
	mux.HandleFunc("/api/predictions", h.HandlePredictions)
	mux.HandleFunc("/api/predictions/bulk", h.HandleBulkPredictions)
	mux.HandleFunc("/api/leaderboard", h.HandleLeaderboard)
	mux.HandleFunc("/api/stats", h.HandleUserStats)

	// Apply middleware
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
