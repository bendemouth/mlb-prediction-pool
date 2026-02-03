package database

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/bendemouth/mlb-prediction-pool/internal/models"
)

// CreatePrediction stores a new prediction
func (db *DB) CreatePrediction(ctx context.Context, prediction *models.Prediction) error {
	item, err := attributevalue.MarshalMap(prediction)
	if err != nil {
		return fmt.Errorf("failed to marshal prediction: %w", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(db.predictionsTable),
		Item:      item,
	}

	_, err = db.client.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to create prediction: %w", err)
	}

	return nil
}

// GetUserPredictions retrieves all predictions for a user
func (db *DB) GetUserPredictions(ctx context.Context, userID string) ([]models.Prediction, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(db.predictionsTable),
		KeyConditionExpression: aws.String("userId = :userId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userId": &types.AttributeValueMemberN{Value: fmt.Sprintf("%s", userID)},
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

// GetPrediction retrieves a specific prediction
func (db *DB) GetPrediction(ctx context.Context, userID string, gameID string) (*models.Prediction, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(db.predictionsTable),
		Key: map[string]types.AttributeValue{
			"userId": &types.AttributeValueMemberS{Value: userID},
			"gameId": &types.AttributeValueMemberS{Value: gameID},
		},
	}

	result, err := db.client.GetItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get prediction: %w", err)
	}

	if result.Item == nil {
		return nil, fmt.Errorf("prediction not found")
	}

	var prediction models.Prediction
	err = attributevalue.UnmarshalMap(result.Item, &prediction)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal prediction: %w", err)
	}

	return &prediction, nil
}

// UpdatePredictionStatus updates the status of a prediction after game completion
func (db *DB) UpdatePredictionStatus(ctx context.Context, userID string, gameID, status, actualWinner string) error {
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(db.predictionsTable),
		Key: map[string]types.AttributeValue{
			"userId": &types.AttributeValueMemberS{Value: userID},
			"gameId": &types.AttributeValueMemberS{Value: gameID},
		},
		UpdateExpression: aws.String("SET #status = :status, actualWinner = :actualWinner"),
		ExpressionAttributeNames: map[string]string{
			"#status": "status", // 'status' is a reserved word in DynamoDB
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":status":       &types.AttributeValueMemberS{Value: status},
			":actualWinner": &types.AttributeValueMemberS{Value: actualWinner},
		},
	}

	_, err := db.client.UpdateItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update prediction: %w", err)
	}

	return nil
}

// GetPredictionsByGame retrieves all predictions for a specific game
func (db *DB) GetPredictionsByGame(ctx context.Context, gameID string) ([]models.Prediction, error) {
	// This requires a GSI on gameId
	// For now, use scan with filter (inefficient but works)

	input := &dynamodb.ScanInput{
		TableName:        aws.String(db.predictionsTable),
		FilterExpression: aws.String("gameId = :gameId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":gameId": &types.AttributeValueMemberS{Value: gameID},
		},
	}

	result, err := db.client.Scan(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to scan predictions: %w", err)
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

// BatchCreatePredictions creates multiple predictions efficiently
func (db *DB) BatchCreatePredictions(ctx context.Context, predictions []models.Prediction) error {
	// DynamoDB BatchWriteItem supports up to 25 items
	const batchSize = 25

	for i := 0; i < len(predictions); i += batchSize {
		end := i + batchSize
		if end > len(predictions) {
			end = len(predictions)
		}

		batch := predictions[i:end]
		writeRequests := make([]types.WriteRequest, 0, len(batch))

		for _, pred := range batch {
			item, err := attributevalue.MarshalMap(pred)
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
