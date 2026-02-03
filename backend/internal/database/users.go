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

	entity := ToUserEntity(user)

	item, err := attributevalue.MarshalMap(entity)
	if err != nil {
		return fmt.Errorf("Failed to marshal user entity: %w", err)
	}

	input := &dynamodb.PutItemInput{
		TableName:           aws.String(db.usersTable),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(PK)"),
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
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userId)},
			"SK": &types.AttributeValueMemberS{Value: "PROFILE"},
		},
	}

	result, err := db.client.GetItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("Failed to get item from DynamoDB: %w", err)
	}

	var entity UserEntity
	err = attributevalue.UnmarshalMap(result.Item, &entity)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal user entity: %w", err)
	}

	return FromUserEntity(&entity), nil
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
