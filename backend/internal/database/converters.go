package database

import (
	"fmt"

	"github.com/bendemouth/mlb-prediction-pool/internal/models"
)

func ToUserEntity(user *models.User) *UserEntity {
	return &UserEntity{
		PK:         fmt.Sprintf("USER#%s", user.Id),
		SK:         "PROFILE",
		UserId:     user.Id,
		Username:   user.Username,
		Email:      user.Email,
		CreatedAt:  user.CreatedAt,
		EntityType: "USER",
	}
}

func FromUserEntity(entity *UserEntity) *models.User {
	return &models.User{
		Id:        entity.UserId,
		Username:  entity.Username,
		Email:     entity.Email,
		CreatedAt: entity.CreatedAt,
	}
}

// ToPredictionEntity converts domain Prediction to database PredictionEntity
func ToPredictionEntity(pred *models.Prediction) *PredictionEntity {
	entity := &PredictionEntity{
		PK:                  fmt.Sprintf("USER#%s", pred.UserId),
		SK:                  fmt.Sprintf("PREDICTION#%s", pred.GameId),
		UserId:              pred.UserId,
		GameId:              pred.GameId,
		HomeScorePredicted:  pred.HomeScorePredicted,
		AwayScorePredicted:  pred.AwayScorePredicted,
		TotalScorePredicted: pred.TotalScorePredicted,
		Confidence:          pred.Confidence,
		PredictedWinnerId:   pred.PredictedWinnerId,
		SubmittedAt:         pred.SubmittedAt,
		EntityType:          "PREDICTION",

		// GSI keys for reverse lookup
		GSI1PK: fmt.Sprintf("GAME#%s", pred.GameId),
		GSI1SK: fmt.Sprintf("USER#%s", pred.UserId),
	}

	if pred.ActualWinnerId != "" {
		entity.ActualWinnerId = &pred.ActualWinnerId
	}
	if pred.WinnerCorrect != nil {
		entity.WinnerCorrect = pred.WinnerCorrect
	}

	return entity
}

// FromPredictionEntity converts database PredictionEntity to domain Prediction
func FromPredictionEntity(entity *PredictionEntity) *models.Prediction {
	pred := &models.Prediction{
		UserId:              entity.UserId,
		GameId:              entity.GameId,
		HomeScorePredicted:  entity.HomeScorePredicted,
		AwayScorePredicted:  entity.AwayScorePredicted,
		TotalScorePredicted: entity.TotalScorePredicted,
		Confidence:          entity.Confidence,
		PredictedWinnerId:   entity.PredictedWinnerId,
		SubmittedAt:         entity.SubmittedAt,
		WinnerCorrect:       entity.WinnerCorrect,
	}

	if entity.ActualWinnerId != nil {
		pred.ActualWinnerId = *entity.ActualWinnerId
	}

	return pred
}
