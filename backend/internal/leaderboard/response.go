package leaderboard

import dbgen "github.com/anssuy/code-colosseum/backend/internal/db/generated"

type EntryResponse struct {
	Rank     int    `json:"rank"`
	Username string `json:"username"`
	Rating   int32  `json:"rating"`
	Wins     int32  `json:"wins"`
	Losses   int32  `json:"losses"`
}

func EntryResponseFrom(user dbgen.ListLeaderboardRow, rank int) EntryResponse {
	return EntryResponse{
		Rank:     rank,
		Username: user.Username,
		Rating:   user.Rating,
		Wins:     user.Wins,
		Losses:   user.Losses,
	}
}
