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

// CompleteGame marks a game as completed and updates related predictions
func (db *DB) CompleteGame(ctx context.Context, gameId string, homeScore int, awayScore int, winnerId string) error {
	// Update game record
	updateGameInput := &dynamodb.UpdateItemInput{
		TableName: aws.String(db.gamesTable),
		Key: map[string]types.AttributeValue{
			"gameId": &types.AttributeValueMemberS{Value: gameId},
		},
		UpdateExpression: aws.String(
			"SET #status = :status, homeScore = :homeScore, awayScore = :awayScore, winner = :winner",
		),
		ExpressionAttributeNames: map[string]string{
			"#status": "status",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":status":    &types.AttributeValueMemberS{Value: "completed"},
			":homeScore": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", homeScore)},
			":awayScore": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", awayScore)},
			":winner":    &types.AttributeValueMemberS{Value: winnerId},
		},
	}

	_, err := db.client.UpdateItem(ctx, updateGameInput)
	if err != nil {
		return fmt.Errorf("failed to update game: %w", err)
	}

	// Get all predictions for the game
	predictions, err := db.GetPredictionsByGame(ctx, gameId)
	if err != nil {
		return fmt.Errorf("failed to get predictions: %w", err)
	}

	// Update each prediction based on the game result
	for _, prediction := range predictions {
		if err := db.updatePredictionsWithResult(ctx, prediction.UserId, gameId, winnerId, homeScore, awayScore); err != nil {
			return fmt.Errorf("failed to update prediction for user %s: %w", prediction.UserId, err)
		}
	}

	return nil
}

func (db *DB) updatePredictionsWithResult(ctx context.Context, userId, gameId, winnerId string, homeScore, awayScore int) error {
	pred, err := db.GetPredictionByUser(ctx, userId, gameId)
	if err != nil {
		return fmt.Errorf("failed to get prediction: %w", err)
	}

	homeScoreError := abs(pred.HomeScorePredicted - float32(homeScore))
	awayScoreError := abs(pred.AwayScorePredicted - float32(awayScore))
	totalScoreError := abs(pred.TotalScorePredicted - float32(homeScore+awayScore))
	winnerCorrect := pred.PredictedWinnerId == winnerId

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(db.predictionsTable),
		Key: map[string]types.AttributeValue{
			"userId": &types.AttributeValueMemberS{Value: userId},
			"gameId": &types.AttributeValueMemberS{Value: gameId},
		},
		UpdateExpression: aws.String(
			"SET actualWinnerId = :actualWinnerId, " +
				"winnerCorrect = :winnerCorrect, " +
				"homeScoreError = :homeScoreError, " +
				"awayScoreError = :awayScoreError, " +
				"totalScoreError = :totalScoreError, ",
		),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":actualWinnerId":  &types.AttributeValueMemberS{Value: winnerId},
			":winnerCorrect":   &types.AttributeValueMemberBOOL{Value: winnerCorrect},
			":homeScoreError":  &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", homeScoreError)},
			":awayScoreError":  &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", awayScoreError)},
			":totalScoreError": &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", totalScoreError)},
		},
	}

	_, err = db.client.UpdateItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update prediction: %w", err)
	}

	return nil
}

func abs(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}
