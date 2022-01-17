package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	firstName    *string            `json: "firstName" validate:"required, min=2, max=100"`
	lastNmae     *string            `json: "firstName" validate:"required, min=2, max=100"`
	password     *string            `json: "password" validate:"required, min=6,"`
	email        *string            `json: "email" validate:"required, email"`
	phone        *string            `json: "phone" validate:"required, min=6,"`
	token        *string            `json: "token"`
	userType     *string            `json: "userType" validate:"required, eq=Admin|eq=User"`
	refreshToken *string            `json: "refreshToken""`
	createdAt    time.Time            `json: "createdAt"`
	updatedAt    time.Time            `json: "updatedAt"`
	userID       *string            `json: "userID"`
}