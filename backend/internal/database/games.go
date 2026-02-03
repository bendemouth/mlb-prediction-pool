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

// CreateGame stores a new game
func (db *DB) CreateGame(ctx context.Context, game *models.Game) error {
	item, err := attributevalue.MarshalMap(game)
	if err != nil {
		return fmt.Errorf("failed to marshal game: %w", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(db.gamesTable),
		Item:      item,
	}

	_, err = db.client.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to create game: %w", err)
	}

	return nil
}

// GetGame retrieves a game by ID
func (db *DB) GetGame(ctx context.Context, gameID string) (*models.Game, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(db.gamesTable),
		Key: map[string]types.AttributeValue{
			"gameId": &types.AttributeValueMemberS{Value: gameID},
		},
	}

	result, err := db.client.GetItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get game: %w", err)
	}

	if result.Item == nil {
		return nil, fmt.Errorf("game not found")
	}

	var game models.Game
	err = attributevalue.UnmarshalMap(result.Item, &game)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal game: %w", err)
	}

	return &game, nil
}

// GetUpcomingGames retrieves games with status "upcoming"
func (db *DB) GetUpcomingGames(ctx context.Context) ([]models.Game, error) {
	input := &dynamodb.ScanInput{
		TableName:        aws.String(db.gamesTable),
		FilterExpression: aws.String("#status = :status"),
		ExpressionAttributeNames: map[string]string{
			"#status": "status",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":status": &types.AttributeValueMemberS{Value: "upcoming"},
		},
	}

	result, err := db.client.Scan(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to scan games: %w", err)
	}

	games := make([]models.Game, 0, len(result.Items))
	for _, item := range result.Items {
		var game models.Game
		if err := attributevalue.UnmarshalMap(item, &game); err != nil {
			return nil, fmt.Errorf("failed to unmarshal game: %w", err)
		}
		games = append(games, game)
	}

	return games, nil
}

// UpdateGameResult updates a game's result
func (db *DB) UpdateGameResult(ctx context.Context, gameID, winner string) error {
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(db.gamesTable),
		Key: map[string]types.AttributeValue{
			"gameId": &types.AttributeValueMemberS{Value: gameID},
		},
		UpdateExpression: aws.String("SET #status = :status, winner = :winner"),
		ExpressionAttributeNames: map[string]string{
			"#status": "status",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":status": &types.AttributeValueMemberS{Value: "completed"},
			":winner": &types.AttributeValueMemberS{Value: winner},
		},
	}

	_, err := db.client.UpdateItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update game: %w", err)
	}

	return nil
}
