package models

import "time"

type User struct {
	Id        string    `json:"id" dynamodbav:"userId"`
	Username  string    `json:"username" dynamodbav:"username"`
	Email     string    `json:"email" dynamodbav:"email"`
	CreatedAt time.Time `json:"created_at" dynamodbav:"createdAt"`
}
