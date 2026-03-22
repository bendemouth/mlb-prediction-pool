package models

import "time"

type ModelMetadata struct {
	ModelId   string    `json:"model_id" dynamodbav:"modelId"`
	ModelName string    `json:"model_name" dynamodbav:"modelName"`
	UserId    string    `json:"user_id" dynamodbav:"userId"`
	FileName  string    `json:"file_name" dynamodbav:"fileName"`
	S3Key     string    `json:"s3_key" dynamodbav:"s3Key"`
	Status    string    `json:"status" dynamodbav:"status"`
	CreatedAt time.Time `json:"created_at" dynamodbav:"createdAt"`
	UpdatedAt time.Time `json:"updated_at" dynamodbav:"updatedAt"`
}
