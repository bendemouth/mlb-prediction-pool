package models

import "time"

type User struct {
	Id        string    `json:"id" dynamodb:"userId"`
	Username  string    `json:"username" dynamodb:"username"`
	Email     string    `json:"email" dynamodb:"email"`
	CreatedAt time.Time `json:"created_at" dynamodb:"createdAt"`
}
