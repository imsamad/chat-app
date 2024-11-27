package handlers

import (
	"chat-app/mongorm"
	"chat-app/utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h *Handler) GetInitialState(w http.ResponseWriter, r *http.Request) *utils.Response {

	userID := r.Context().Value("user").(primitive.ObjectID)

	friendshipCollection := h.DB.Collection("friendships")

	pipeline := mongo.Pipeline{
		{
			{Key: "$match", Value: bson.D{
				{Key: "$or", Value: bson.A{
					bson.M{"user_a": userID},
					bson.M{"user_b": userID},
				}},
			}},
		},
		{
			{Key: "$addFields", Value: bson.D{
				{Key: "friend", Value: bson.D{
					{Key: "$cond", Value: bson.D{
						{Key: "if", Value: bson.D{
							{Key: "$eq", Value: bson.A{"$user_a", userID}},
						}},
						{Key: "then", Value: "$user_b"},
						{Key: "else", Value: "$user_a"},
					}},
				}},
			}},
		},
		{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "users"},
				{Key: "localField", Value: "friend"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "friend"},
			}},
		},
		{
			{Key: "$unwind", Value: bson.D{
				{Key: "path", Value: "$friend"},
				{Key: "preserveNullAndEmptyArrays", Value: true},
			}},
		},
		{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "messages"},
				{Key: "localField", Value: "_id"},
				{Key: "foreignField", Value: "friendship_id"},
				{Key: "as", Value: "messages"},
			}},
		},
		{
			{Key: "$addFields", Value: bson.D{
				{Key: "messages", Value: bson.D{
					{Key: "$sortArray", Value: bson.D{
						{Key: "input", Value: "$messages"},
						{Key: "sortBy", Value: bson.D{
							{Key: "last_id", Value: -1},
						}},
					}},
				}},
			}},
		},
		{{Key: "$unset", Value: bson.A{"user_a", "user_b"}}},
	}

	ctx := context.Background()
	cursor, err := friendshipCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return &utils.Response{
			Code:    http.StatusConflict,
			Message: "Data not found",
		}
	}
	defer cursor.Close(ctx)

	type Response struct {
		mongorm.Model `bson:",inline" json:",inline"`
		Friend        mongorm.UserModel      `bson:"friend" json:"friend"`
		Messages      []mongorm.MessageModel `bson:"messages" json:"messages"`
	}

	// Decode the results
	var results []Response
	if err := cursor.All(ctx, &results); err != nil {
		return &utils.Response{
			Code:    http.StatusConflict,
			Message: "Data not found",
		}
	}

	return &utils.Response{
		Code:    http.StatusOK,
		Message: results,
	}

}

func (h *Handler) CreateMessage(w http.ResponseWriter, r *http.Request) *utils.Response {
	type Body struct {
		Text       string             `json:"text"`
		Friendship primitive.ObjectID `json:"friendship"`
		// From primitive.ObjectID `json:"from"`
	}

	var body Body
	json.NewDecoder(r.Body).Decode(&body)
	user_id := r.Context().Value("user").(primitive.ObjectID)

	var last_message mongorm.MessageModel
	err := last_message.Read(context.Background(), h.DB, "messages", bson.M{"friendship_id": body.Friendship}, &last_message, bson.D{{Key: "last_id", Value: -1}})

	json.NewEncoder(os.Stdout).Encode(last_message)
	if err != nil {
		fmt.Println("error while fetching last message doc: ", err)
		return &utils.Response{
			Code:    http.StatusConflict,
			Message: "not found last doc",
		}
	}

	var last_id int64 = 0

	if last_message.LastId != 0 {
		last_id = last_message.LastId
	}

	message := mongorm.MessageModel{
		From:         user_id,
		Text:         body.Text,
		FriendshipId: body.Friendship,
		LastId:       last_id + 1,
	}

	err = message.Create(context.Background(), h.DB, "messages", &message)

	if err != nil {
		return &utils.Response{
			Code:    http.StatusOK,
			Message: "Please try again!",
		}
	}

	return &utils.Response{
		Code: http.StatusOK,
		Message: map[string]interface{}{
			"message": message,
		},
	}
}
