package database

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DB struct {
	client           *dynamodb.Client
	usersTable       string
	predictionsTable string
	leaderboardTable string
	gamesTable       string
	modelsTable      string
}

type DBConfig struct {
	Region           string
	Endpoint         string
	UsersTable       string
	PredictionsTable string
	LeaderboardTable string
	GamesTable       string
	ModelsTable      string
}

// NewDB creates a new database connection
func NewDB(ctx context.Context, cfg DBConfig) (*DB, error) {
	// Load AWS configuration
	awsCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(cfg.Region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create DynamoDB client
	var client *dynamodb.Client
	if cfg.Endpoint != "" {
		// Use local DynamoDB (for development)
		client = dynamodb.NewFromConfig(awsCfg, func(o *dynamodb.Options) {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
		})
	} else {
		client = dynamodb.New(dynamodb.Options{Credentials: awsCfg.Credentials, Region: awsCfg.Region})
	}

	db := &DB{
		client:           client,
		usersTable:       cfg.UsersTable,
		predictionsTable: cfg.PredictionsTable,
		leaderboardTable: cfg.LeaderboardTable,
		gamesTable:       cfg.GamesTable,
		modelsTable:      cfg.ModelsTable,
	}

	return db, nil
}

func NewDBFromEnv(ctx context.Context) (*DB, error) {
	cfg := DBConfig{
		Region:           getEnv("DYNAMODB_REGION", "us-west-2"),
		Endpoint:         os.Getenv("DYNAMODB_ENDPOINT"), // Optional
		UsersTable:       getEnv("DYNAMODB_USERS_TABLE", "mlb-prediction-pool-users"),
		PredictionsTable: getEnv("DYNAMODB_PREDICTIONS_TABLE", "mlb-prediction-pool-predictions"),
		LeaderboardTable: getEnv("DYNAMODB_LEADERBOARD_TABLE", "mlb-prediction-pool-leaderboard"),
		GamesTable:       getEnv("DYNAMODB_GAMES_TABLE", "mlb-prediction-pool-games"),
		ModelsTable:      getEnv("DYNAMODB_MODELS_TABLE", "mlb-prediction-pool-models"),
	}

	return NewDB(ctx, cfg)
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func (db *DB) Close() error {
	// DynamoDB client does not require explicit closure
	return nil
}
