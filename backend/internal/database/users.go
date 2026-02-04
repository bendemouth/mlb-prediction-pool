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

// CreateUser adds a new user to the Users table
func (db *DB) CreateUser(ctx context.Context, user *models.User) error {
	user.CreatedAt = time.Now()

	item, err := attributevalue.MarshalMap(user)
	if err != nil {
		return fmt.Errorf("Failed to marshal user entity: %w", err)
	}

	input := &dynamodb.PutItemInput{
		TableName:           aws.String(db.usersTable),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(userId)"),
	}

	_, err = db.client.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("Failed to put item in DynamoDB: %w", err)
	}

	return nil
}

// GetUser retrieves a user by userId from the Users table
func (db *DB) GetUser(ctx context.Context, userId string) (*models.User, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(db.usersTable),
		Key: map[string]types.AttributeValue{
			"userId": &types.AttributeValueMemberS{Value: userId},
		},
	}

	result, err := db.client.GetItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("Failed to get item from DynamoDB: %w", err)
	}

	var user models.User
	err = attributevalue.UnmarshalMap(result.Item, &user)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal user entity: %w", err)
	}

	return &user, nil
}

// ListUsers retrieves all users from the Users table
func (db *DB) ListUsers(ctx context.Context) ([]*models.User, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(db.usersTable),
	}

	result, err := db.client.Scan(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("Failed to scan DynamoDB: %w", err)
	}

	users := make([]*models.User, 0, len(result.Items))
	for _, item := range result.Items {
		var user models.User
		err = attributevalue.UnmarshalMap(item, &user)
		if err != nil {
			return nil, fmt.Errorf("Failed to unmarshal user: %w", err)
		}
		users = append(users, &user)
	}

	return users, nil
}
