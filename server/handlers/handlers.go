package handlers

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
	DB *mongo.Database
}

func New(db *mongo.Database) *Handler {
	return &Handler{
		DB: db,
	}
}
