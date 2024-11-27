package mongorm

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Model struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type UserModel struct {
	Model    `bson:",inline"`
	Password string `bson:"password" json:"-"`
	Email    string `bson:"email" json:"email"`
	Name     string `bson:"name" json:"name"`
}

type FriendshipModel struct {
	Model `bson:",inline" json:",inline"`
	UserA primitive.ObjectID `bson:"user_a" json:"user_a"`
	UserB primitive.ObjectID `bson:"user_b" json:"user_b"`
}

type MessageModel struct {
	Model        `bson:",inline" json:",inline"`
	Text         string             `bson:"text" json:"text"`
	FriendshipId primitive.ObjectID `bson:"friendship_id" json:"friendship_id"`
	From         primitive.ObjectID `bson:"from" json:"from"`
	LastId       int64              `bson:"last_id" json:"last_id"`
	Sent         time.Time          `bson:"sent" json:"sent"`
	Seen         time.Time          `bson:"seen" json:"seen"`
}
