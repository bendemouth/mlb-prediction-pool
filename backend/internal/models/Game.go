package models

import "time"

type Game struct {
	GameId     string    `json:"game_id"`
	Date       time.Time `json:"date"`
	HomeTeam   string    `json:"home_team"`
	HomeTeamId int       `json:"home_id"`
	AwayTeamId int       `json:"away_id"`
	AwayTeam   string    `json:"away_team"`
	HomeScore  int       `json:"home_score"`
	AwayScore  int       `json:"away_score"`
}
