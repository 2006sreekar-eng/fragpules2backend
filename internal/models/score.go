package models

import "time"

type Score struct {
	ID              string    `json:"id"`
	PlayerName      string    `json:"player_name"`
	Score           int       `json:"score"`
	Accuracy        float64   `json:"accuracy"`
	AvgReactionTime int       `json:"avg_reaction_time"`
	CreatedAt       time.Time `json:"created_at"`
}

type CreateScoreRequest struct {
	PlayerName      string  `json:"player_name"`
	Score           int     `json:"score"`
	Accuracy        float64 `json:"accuracy"`
	AvgReactionTime int     `json:"avg_reaction_time"`
}