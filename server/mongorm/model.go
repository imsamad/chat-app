package mongorm

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Model struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	CreateAt  time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

type UserModel struct {
	Model    `bson:",inline"`
	Password string               `bson:"password"`
	Email    string               `bson:"email"`
	Name     string               `bson:"name"`
	Friends  []primitive.ObjectID `bson:"friends"`
}

type MessageModel struct {
	Model `bson:",inline"`
	Text  string             `bson:"text"`
	From  primitive.ObjectID `bson:"from"`
	To    primitive.ObjectID `bson:"to"`
	SeqId int64              `bson:"seqId"`
}
