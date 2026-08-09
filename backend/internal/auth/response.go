package auth

import (
	dbgen "github.com/anssuy/code-colosseum/backend/internal/db/generated"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserResponse struct {
	ID        pgtype.UUID        `json:"id"`
	Username  string             `json:"username"`
	Email     string             `json:"email"`
	Rating    int32              `json:"rating"`
	Wins      int32              `json:"wins"`
	Losses    int32              `json:"losses"`
	CreatedAt pgtype.Timestamptz `json:"createdAt"`
}

func UserResponseFrom(user dbgen.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Rating:    user.Rating,
		Wins:      user.Wins,
		Losses:    user.Losses,
		CreatedAt: user.CreatedAt,
	}
}
