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
	gamesTable       string
	modelsTable      string
}

type DBConfig struct {
	Region           string
	Endpoint         string
	UsersTable       string
	PredictionsTable string
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
		// Use local DynamoDB for development
		client = dynamodb.NewFromConfig(awsCfg, func(o *dynamodb.Options) {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
		})
	} else {
		client = dynamodb.NewFromConfig(awsCfg)
	}

	db := &DB{
		client:           client,
		usersTable:       cfg.UsersTable,
		predictionsTable: cfg.PredictionsTable,
		gamesTable:       cfg.GamesTable,
		modelsTable:      cfg.ModelsTable,
	}

	return db, nil
}

func NewDBFromEnv(ctx context.Context) (*DB, error) {
	cfg := DBConfig{
		Region:           getEnv("DYNAMODB_REGION", "us-west-2"),
		Endpoint:         getEnv("DYNAMODB_ENDPOINT", ""),
		UsersTable:       getEnv("DYNAMODB_USERS_TABLE", "mlb-prediction-pool-users"),
		PredictionsTable: getEnv("DYNAMODB_PREDICTIONS_TABLE", "mlb-prediction-pool-predictions"),
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

func (db *DB) HealthCheck(ctx context.Context) error {
	_, err := db.client.ListTables(ctx, &dynamodb.ListTablesInput{
		Limit: aws.Int32(1),
	})
	return err
}
