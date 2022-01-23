package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RegisterTournament struct{
	ID						primitive.ObjectID		`bson:"_id" validate:"required"`
	UserName				string					`json:"UserName" validate:"required,min=2,max=100"`
	Created_at				time.Time				`json:"Created_at" validate:"required"`
	Updated_at				time.Time				`json:"Updated_at" validate:"required"`
	TournamentId			string					`json:"TournamentId" validate:"required"`
	RegisterTournamentId	string					`json:"RegisterTournamentId" validate:"required"`
	User_id					string					`json:"User_id" validate:"required"`
}