package handlers

import (
	"chat-app/mongorm"
	"chat-app/utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// allow to add friend - i.e. at the end it create a FriendshipModel document
// it expect to send friend either 1. email or 2. ObjectId of the friend
func (h *Handler) AddFriend(w http.ResponseWriter, r *http.Request) *utils.Response {

	type Body struct {
		Friend string `json:"friend"`
		Email  string `json:"email"`
	}

	var body Body
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return &utils.Response{
			Code: http.StatusConflict,
			Message: map[string]string{
				"message": "Provide friend id or email",
			},
		}
	}

	var friend_id primitive.ObjectID

	if body.Email == "" {
		friend_id, err = primitive.ObjectIDFromHex(body.Friend)
		if err != nil {
			return &utils.Response{
				Message: map[string]string{
					"friend": "It must be valid ObjectId",
				},
			}
		}
	} else {
		var user mongorm.UserModel
		err := user.Read(context.Background(), h.DB, "users", bson.M{"email": body.Email}, &user, nil)

		if err != nil {
			return &utils.Response{
				Code: http.StatusConflict,
				Message: map[string]string{
					"message": "Email not found",
				},
			}
		}
		friend_id = user.ID
	}

	user_id := r.Context().Value("user").(primitive.ObjectID)
	if user_id == friend_id {
		return &utils.Response{
			Code: http.StatusConflict,
			Message: map[string]string{
				"error": "It must be not be your own id",
			},
		}
	}

	friendship := mongorm.FriendshipModel{
		UserA: user_id,
		UserB: friend_id,
	}

	err = friendship.Create(context.Background(), h.DB, "friendships", &friendship)

	if err != nil {
		return &utils.Response{
			Code: http.StatusInternalServerError,
			Message: map[string]string{
				"error": "please try again!",
			},
		}
	}

	return &utils.Response{
		Code: http.StatusOK,
		Message: map[string]string{
			"message": "added  successfully!",
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

	err = mongorm.Delete(context.Background(), h.DB, "friendships", bson.D{
		{Key: "$and", Value: bson.A{
			bson.M{"user_id": user_id},
			bson.M{"friend_id": body.Friend},
		}},
	},
	)

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
			"message": "removed successfully!",
		},
	}
}

func ObjectIDsToHex(ids []primitive.ObjectID) []string {
	result := make([]string, len(ids))

	for i, id := range ids {
		result[i] = id.Hex()
	}

	return result
}

func (h *Handler) GetMyFriends(w http.ResponseWriter, r *http.Request) *utils.Response {
	user_id := r.Context().Value("user").(primitive.ObjectID)

	var friends []mongorm.FriendshipModel
	err := mongorm.ReadAll(context.Background(), h.DB, "friendships", bson.D{{Key: "$or", Value: bson.A{bson.M{"user_a": user_id},
		bson.M{"user_b": user_id}}}}, &friends, bson.M{})

	// fmt.Println(len(friends))
	// for _, friend := range friends {
	// 	fmt.Println(friend)
	// }
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{
			{Key: "$or", Value: bson.A{
				bson.M{"user_a": user_id},
				bson.M{"user_b": user_id},
			}},
		}}},
		{{Key: "$addFields", Value: bson.D{
			{Key: "friend", Value: bson.D{
				{Key: "$cond", Value: bson.D{
					{Key: "if", Value: bson.D{
						{Key: "$eq", Value: bson.A{"$user_a", user_id}},
					}},
					{Key: "then", Value: "$user_b"},
					{Key: "else", Value: "$user_a"},
				}},
			}},
		}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"},
			{Key: "localField", Value: "friend"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "friend"},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$friend"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
		{{Key: "$unset", Value: bson.A{"user_a", "user_b"}}},
	}

	friendshipColl := h.DB.Collection("friendships")
	ctx := context.Background()
	cursor, err := friendshipColl.Aggregate(ctx, pipeline)
	if err != nil {
		fmt.Println(err)
		return &utils.Response{
			Code:    http.StatusConflict,
			Message: "Data not found",
		}
	}
	defer cursor.Close(ctx)
	type Response struct {
		mongorm.FriendshipModel `bson:",inline" json:",inline"`
		Friend                  mongorm.UserModel `bson:"friend" json:"friend"`
		// Friend primitive.ObjectID `bson:"friend" json:"friend"`
		// UserA                   string            `bson:"-" json:"-"`
		// UserB                   string            `bson:"-" json:"-"`
	}
	var results []Response
	if err := cursor.All(ctx, &results); err != nil {
		fmt.Println("err: ", err)
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

/*
[
    {
        "id": "6741c2d43caac9cf474a552d",
        "created_at": "2024-11-23T11:56:04.871Z",
        "updated_at": "2024-11-23T11:56:04.871Z",
        "user_a": "6741c2d43caac9cf474a5523",
        "user_b": "6741c2d43caac9cf474a5527",
        "friend": {
            "id": "6741c2d43caac9cf474a5527",
            "created_at": "2024-11-23T11:56:04.867Z",
            "updated_at": "2024-11-23T11:56:04.867Z",
            "email": "user4@gmail.com",
            "name": "User 4"
        }
    },
    {
        "id": "6741c2d43caac9cf474a552e",
        "created_at": "2024-11-23T11:56:04.871Z",
        "updated_at": "2024-11-23T11:56:04.871Z",
        "user_a": "6741c2d43caac9cf474a5523",
        "user_b": "6741c2d43caac9cf474a5525",
        "friend": {
            "id": "6741c2d43caac9cf474a5525",
            "created_at": "2024-11-23T11:56:04.867Z",
            "updated_at": "2024-11-23T11:56:04.867Z",
            "email": "user2@gmail.com",
            "name": "User 2"
        }
    },
    {
        "id": "6741c2d43caac9cf474a552f",
        "created_at": "2024-11-23T11:56:04.872Z",
        "updated_at": "2024-11-23T11:56:04.872Z",
        "user_a": "6741c2d43caac9cf474a5523",
        "user_b": "6741c2d43caac9cf474a5529",
        "friend": {
            "id": "6741c2d43caac9cf474a5529",
            "created_at": "2024-11-23T11:56:04.868Z",
            "updated_at": "2024-11-23T11:56:04.868Z",
            "email": "user6@gmail.com",
            "name": "User 6"
        }
    },
    {
        "id": "6741c2d43caac9cf474a5530",
        "created_at": "2024-11-23T11:56:04.872Z",
        "updated_at": "2024-11-23T11:56:04.872Z",
        "user_a": "6741c2d43caac9cf474a5523",
        "user_b": "6741c2d43caac9cf474a552b",
        "friend": {
            "id": "6741c2d43caac9cf474a552b",
            "created_at": "2024-11-23T11:56:04.869Z",
            "updated_at": "2024-11-23T11:56:04.869Z",
            "email": "user8@gmail.com",
            "name": "User 8"
        }
    },
    {
        "id": "6741c2d43caac9cf474a5531",
        "created_at": "2024-11-23T11:56:04.872Z",
        "updated_at": "2024-11-23T11:56:04.872Z",
        "user_a": "6741c2d43caac9cf474a5523",
        "user_b": "6741c2d43caac9cf474a5529",
        "friend": {
            "id": "6741c2d43caac9cf474a5529",
            "created_at": "2024-11-23T11:56:04.868Z",
            "updated_at": "2024-11-23T11:56:04.868Z",
            "email": "user6@gmail.com",
            "name": "User 6"
        }
    }
]


*/
