// internal/database/entities.go - NEW FILE
package database

import "time"

// UserEntity represents how User data is stored in DynamoDB
type UserEntity struct {
	PK         string    `dynamodbav:"PK"` // "USER#<userId>"
	SK         string    `dynamodbav:"SK"` // "PROFILE"
	UserId     string    `dynamodbav:"userId"`
	Username   string    `dynamodbav:"username"`
	Email      string    `dynamodbav:"email"`
	CreatedAt  time.Time `dynamodbav:"createdAt"`
	EntityType string    `dynamodbav:"entityType"` // "USER"
}

// PredictionEntity represents how Prediction data is stored in DynamoDB
type PredictionEntity struct {
	PK                  string    `dynamodbav:"PK"` // "USER#<userId>"
	SK                  string    `dynamodbav:"SK"` // "PREDICTION#<gameId>"
	UserId              string    `dynamodbav:"userId"`
	GameId              string    `dynamodbav:"gameId"`
	HomeScorePredicted  float32   `dynamodbav:"homeScorePredicted"`
	AwayScorePredicted  float32   `dynamodbav:"awayScorePredicted"`
	TotalScorePredicted float32   `dynamodbav:"totalScorePredicted"`
	Confidence          float32   `dynamodbav:"confidence"`
	PredictedWinnerId   string    `dynamodbav:"predictedWinnerId"`
	ActualWinnerId      *string   `dynamodbav:"actualWinnerId,omitempty"`
	WinnerCorrect       *bool     `dynamodbav:"winnerCorrect,omitempty"`
	HomeScoreError      *float32  `dynamodbav:"homeScoreError,omitempty"`
	AwayScoreError      *float32  `dynamodbav:"awayScoreError,omitempty"`
	TotalScoreError     *float32  `dynamodbav:"totalScoreError,omitempty"`
	SubmittedAt         time.Time `dynamodbav:"submittedAt"`
	EntityType          string    `dynamodbav:"entityType"` // "PREDICTION"

	// GSI-1 keys for querying predictions by game
	GSI1PK string `dynamodbav:"GSI1PK"` // "GAME#<gameId>"
	GSI1SK string `dynamodbav:"GSI1SK"` // "USER#<userId>"
}
