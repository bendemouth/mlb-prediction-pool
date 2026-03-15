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

// CreateModel adds a new model to the Models table
func (db *DB) CreateModel(ctx context.Context, model *models.ModelMetadata) error {
	model.CreatedAt = time.Now()
	model.UpdatedAt = time.Now()

	item, err := attributevalue.MarshalMap(model)
	if err != nil {
		return fmt.Errorf("failed to marshal model entity: %w", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(db.modelsTable),
		Item:      item,
	}

	_, err = db.client.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to put model item in DynamoDB: %w", err)
	}

	return nil
}

// GetModelsByUserId retrieves all models for a specific user
func (db *DB) GetModelsByUserId(ctx context.Context, userId string) ([]*models.ModelMetadata, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(db.modelsTable),
		KeyConditionExpression: aws.String("userId = :userId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userId": &types.AttributeValueMemberS{Value: userId},
		},
	}

	result, err := db.client.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query models from DynamoDB: %w", err)
	}

	modelList := make([]*models.ModelMetadata, 0, len(result.Items))
	for _, item := range result.Items {
		var model models.ModelMetadata
		err = attributevalue.UnmarshalMap(item, &model)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal model: %w", err)
		}
		modelList = append(modelList, &model)
	}

	return modelList, nil
}

// GetModelById retrieves a specific model by ID
func (db *DB) GetModelById(ctx context.Context, modelId string, userId string) (*models.ModelMetadata, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(db.modelsTable),
		Key: map[string]types.AttributeValue{
			"userId":  &types.AttributeValueMemberS{Value: userId},
			"modelId": &types.AttributeValueMemberS{Value: modelId},
		},
	}

	result, err := db.client.GetItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get model from DynamoDB: %w", err)
	}

	if result.Item == nil {
		return nil, fmt.Errorf("model not found")
	}

	var model models.ModelMetadata
	err = attributevalue.UnmarshalMap(result.Item, &model)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal model: %w", err)
	}

	return &model, nil
}

// DeleteModel removes a model from the Models table
func (db *DB) DeleteModel(ctx context.Context, modelId string, userId string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(db.modelsTable),
		Key: map[string]types.AttributeValue{
			"userId":  &types.AttributeValueMemberS{Value: userId},
			"modelId": &types.AttributeValueMemberS{Value: modelId},
		},
	}

	_, err := db.client.DeleteItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete model from DynamoDB: %w", err)
	}

	return nil
}

// UpdateModelStatus updates the status of a model
func (db *DB) UpdateModelStatus(ctx context.Context, modelId string, userId string, status string) error {
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(db.modelsTable),
		Key: map[string]types.AttributeValue{
			"userId":  &types.AttributeValueMemberS{Value: userId},
			"modelId": &types.AttributeValueMemberS{Value: modelId},
		},
		UpdateExpression: aws.String("SET #status = :status, updatedAt = :updatedAt"),
		ExpressionAttributeNames: map[string]string{
			"#status": "status",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":status":    &types.AttributeValueMemberS{Value: status},
			":updatedAt": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", time.Now().UnixMilli())},
		},
	}

	_, err := db.client.UpdateItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update model status in DynamoDB: %w", err)
	}

	return nil
}
