package models

import "time"

type ModelMetadata struct {
	ModelId   string    `json:"model_id"`
	ModelName string    `json:"model_name"`
	UserId    string    `json:"user_id"`
	FileName  string    `json:"file_name"`
	S3Key     string    `json:"s3_key"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
