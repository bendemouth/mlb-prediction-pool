package services

import (
	"context"

	"github.com/bendemouth/mlb-prediction-pool/internal/database"
)

type HealthcheckService struct {
	db *database.DB
}

func NewHealthcheckService(db *database.DB) *HealthcheckService {
	return &HealthcheckService{db: db}
}

func (service *HealthcheckService) HealthCheck(ctx context.Context) map[string]string {
	status := map[string]string{
		"service":  "healthy",
		"database": "unknown",
	}

	if err := service.db.HealthCheck(ctx); err != nil {
		status["database"] = "unhealthy"
		return status
	}
	status["database"] = "healthy"
	return status
}
