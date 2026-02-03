package models

import "time"

type Game struct {
	GameId     string    `json:"game_id" dynamodbav:"gameId"`
	Date       time.Time `json:"date" dynamodbav:"date"`
	HomeTeam   string    `json:"home_team" dynamodbav:"homeTeam"`
	HomeTeamId int       `json:"home_id" dynamodbav:"homeTeamId"`
	AwayTeamId int       `json:"away_id" dynamodbav:"awayTeamId"`
	AwayTeam   string    `json:"away_team" dynamodbav:"awayTeam"`
	HomeScore  int       `json:"home_score" dynamodbav:"homeScore"`
	AwayScore  int       `json:"away_score" dynamodbav:"awayScore"`
	Status     string    `json:"status" dynamodbav:"status"`
	Winner     string    `json:"winner,omitempty" dynamodbav:"winner,omitempty"`
}
