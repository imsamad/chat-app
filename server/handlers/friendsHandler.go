package handlers

import (
	"chat-app/mongorm"
	"chat-app/utils"
	"context"
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Handler) AddFriend(w http.ResponseWriter, r *http.Request) *utils.Response {

	type Body struct {
		Friend string `json:"friend"`
	}

	var body Body
	json.NewDecoder(r.Body).Decode(&body)

	friend_id, err := primitive.ObjectIDFromHex(body.Friend)

	if err != nil {
		return &utils.Response{
			Message: map[string]string{
				"friend": "It must be valid ObjectId",
			},
		}
	}

	user_id := r.Context().Value("user").(primitive.ObjectID)

	if user_id == friend_id {
		return &utils.Response{
			Code: http.StatusConflict,
			Message: map[string]string{
				"friend": "It must be not be your own id",
			},
		}
	}
	var user mongorm.UserModel

	// err = user.Update(context.Background(), h.DB, "user", bson.M{"_id": user_id}, bson.D{{"$push", bson.D{{"friends", friend_id}}}})
	err = user.Update(
		context.Background(),
		h.DB,
		"users",
		bson.M{"_id": user_id},
		bson.D{{Key: "$addToSet", Value: bson.M{"friends": friend_id}}},
	)

	if err != nil {
		return &utils.Response{
			Code: http.StatusInternalServerError,
			Message: map[string]string{
				"message": "unable to update the record, please try again!",
			},
		}
	}

	return &utils.Response{
		Code: http.StatusOK,
		Message: map[string]string{
			"friend": "added  successfully!",
		},
	}
}

func (h *Handler) RemoveFriend(w http.ResponseWriter, r *http.Request) *utils.Response {
	type Body struct {
		Friend primitive.ObjectID `json:"friend"`
	}

	var body Body
	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		return &utils.Response{
			Code: http.StatusConflict,
			Message: map[string]string{
				"friend": "it must be valid friend id",
			},
		}
	}

	user_id := r.Context().Value("user").(primitive.ObjectID)

	var user mongorm.UserModel

	err = user.Update(context.Background(), h.DB, "users", bson.M{"_id": user_id}, bson.D{{Key: "$pull", Value: bson.M{"friends": body.Friend}}})

	if err != nil {
		return &utils.Response{
			Code: http.StatusConflict,
			Message: map[string]string{
				"message": "Please try again!",
			},
		}
	}

	return &utils.Response{
		Code: http.StatusOK,
		Message: map[string]string{
			"friend": "removed  successfully!",
		},
	}

}

func (h *Handler) GetMyFriends(w http.ResponseWriter, r *http.Request) *utils.Response {
	user_id := r.Context().Value("user").(primitive.ObjectID)

	var user mongorm.UserModel

	err := user.Read(context.Background(), h.DB, "users", bson.M{"_id": user_id}, &user)

	if err != nil {
		return &utils.Response{
			Code: http.StatusNoContent,
			Message: map[string]string{
				"friends": "Record not found",
			},
		}
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"friends": user.Friends,
	})

	return nil
}
