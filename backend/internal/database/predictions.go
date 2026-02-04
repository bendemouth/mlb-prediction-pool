package database

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/bendemouth/mlb-prediction-pool/internal/models"
)

// CreatePrediction stores a new prediction
func (db *DB) CreatePrediction(ctx context.Context, prediction *models.Prediction) error {
	prediction.SubmittedAt = time.Now()

	item, err := attributevalue.MarshalMap(prediction)
	if err != nil {
		return fmt.Errorf("failed to marshal prediction: %w", err)
	}

	input := dynamodb.PutItemInput{
		TableName: aws.String(db.predictionsTable),
		Item:      item,
	}

	_, err = db.client.PutItem(ctx, &input)
	if err != nil {
		return fmt.Errorf("failed to create prediction: %w", err)
	}

	return nil
}

// GetUserPredictions retrieves all predictions for a user
func (db *DB) GetUserPredictions(ctx context.Context, userId string) ([]models.Prediction, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(db.predictionsTable),
		KeyConditionExpression: aws.String("userId = :userId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userId": &types.AttributeValueMemberS{Value: userId},
		},
	}

	result, err := db.client.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query predictions: %w", err)
	}

	predictions := make([]models.Prediction, 0, len(result.Items))

	for _, item := range result.Items {
		var prediction models.Prediction
		if err := attributevalue.UnmarshalMap(item, &prediction); err != nil {
			return nil, fmt.Errorf("failed to unmarshal prediction: %w", err)
		}
		predictions = append(predictions, prediction)
	}

	return predictions, nil
}

// GetPredictionByUser retrieves a specific prediction by userId and gameId
func (db *DB) GetPredictionByUser(ctx context.Context, userId, gameId string) (*models.Prediction, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(db.predictionsTable),
		Key: map[string]types.AttributeValue{
			"userId": &types.AttributeValueMemberS{Value: userId},
			"gameId": &types.AttributeValueMemberS{Value: gameId},
		},
	}

	result, err := db.client.GetItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get prediction: %w", err)
	}

	if result.Item == nil {
		return nil, nil // Prediction not found
	}

	var prediction models.Prediction
	if err := attributevalue.UnmarshalMap(result.Item, &prediction); err != nil {
		return nil, fmt.Errorf("failed to unmarshal prediction: %w", err)
	}

	return &prediction, nil
}

// GetPredictionsByGame retrieves all predictions for a specific game
func (db *DB) GetPredictionsByGame(ctx context.Context, gameId string) ([]models.Prediction, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(db.predictionsTable),
		IndexName:              aws.String("GameIdIndex"), // Your GSI name
		KeyConditionExpression: aws.String("gameId = :gameId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":gameId": &types.AttributeValueMemberS{Value: gameId},
		},
	}

	result, err := db.client.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query predictions: %w", err)
	}

	predictions := make([]models.Prediction, 0, len(result.Items))
	for _, item := range result.Items {
		var prediction models.Prediction
		if err := attributevalue.UnmarshalMap(item, &prediction); err != nil {
			return nil, fmt.Errorf("failed to unmarshal prediction: %w", err)
		}
		predictions = append(predictions, prediction)
	}

	return predictions, nil
}

// BatchCreatePredictions creates multiple predictions in a single batch operation
func (db *DB) BatchCreatePredictions(ctx context.Context, predictions []models.Prediction) error {
	const batchSize = 25 // DynamoDB batch write limit
	now := time.Now()

	for i := 0; i < len(predictions); i += batchSize {
		end := i + batchSize
		// Adjust end if it exceeds the slice length
		if end > len(predictions) {
			end = len(predictions)
		}

		batch := predictions[i:end]
		writeRequests := make([]types.WriteRequest, 0, len(batch))

		for _, prediction := range batch {
			prediction.SubmittedAt = now

			item, err := attributevalue.MarshalMap(prediction)
			if err != nil {
				return fmt.Errorf("failed to marshal prediction: %w", err)
			}

			writeRequests = append(writeRequests, types.WriteRequest{
				PutRequest: &types.PutRequest{
					Item: item,
				},
			})
		}

		input := &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				db.predictionsTable: writeRequests,
			},
		}

		_, err := db.client.BatchWriteItem(ctx, input)
		if err != nil {
			return fmt.Errorf("failed to batch write predictions: %w", err)
		}
	}

	return nil
}
