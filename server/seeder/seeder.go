package seeder

import (
	"chat-app/mongorm"
	"chat-app/utils"
	"context"
	"math/rand/v2"
	"time"

	"fmt"
	"os"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
)

var USER_COUNT = 10
var FRIENDS_PER_USER = 6
var MESSAGES_COUNTS = 30

func RunSeed() {
	pwd, err := utils.HashPassword("123456")
	if err != nil {
		fmt.Println("pwd not hased")
		os.Exit(1)
	}

	mongoClient, err := mongorm.Connect(os.Getenv("MONGODB_URI"))

	if err != nil {
		fmt.Println("error while connecting to mongo: \n reason ", err)
		os.Exit(1)
	}
	db := mongoClient.Database("chat-db")

	if err != nil {
		fmt.Println("pwd not hased")
		os.Exit(1)
	}

	for i := 0; i < USER_COUNT; i++ {
		user := mongorm.UserModel{
			Name:     fmt.Sprintf("User %s", strconv.Itoa(i)),
			Email:    fmt.Sprintf("user%s@gmail.com", strconv.Itoa(i)),
			Password: pwd,
		}
		err := user.Create(context.Background(), db, "users", &user)
		if err != nil {
			fmt.Println("error while saving: ", user)
		}
	}

	var users []mongorm.UserModel
	err = mongorm.ReadAll(context.Background(), db, "users", bson.M{}, &users, nil)
	fmt.Println(len(users), "users seeded!")

	if err != nil {
		fmt.Println("error while fetching all users: ", err)
	}

	for i, user := range users {

		for j := 0; j < FRIENDS_PER_USER; j++ {
			randNo := rand.IntN(len(users))
			for randNo == 0 || i == randNo || randNo >= len(users) {
				randNo = rand.IntN(len(users))
			}

			existingFriendship := mongorm.FriendshipModel{}

			filter := bson.M{
				"$or": []bson.M{
					{"user_a": user.ID, "user_b": users[randNo].ID},
					{"user_a": users[randNo].ID, "user_b": user.ID},
				},
			}

			err := existingFriendship.Read(context.Background(), db, "friendships", filter, &existingFriendship, nil)
			if err == nil {
				// Friendship already exists, skip creating
				fmt.Printf("Friendship already exists between %s and %s\n", user.ID, users[randNo].ID)
				continue
			}

			friendship := mongorm.FriendshipModel{
				UserA: user.ID,
				UserB: users[randNo].ID,
			}

			err = friendship.Create(context.Background(), db, "friendships", &friendship)

			if err != nil {
				fmt.Println("error while creating friendship doc for: ", friendship)
			}
		}
	}

	var friendships []mongorm.FriendshipModel

	err = mongorm.ReadAll(context.Background(), db, "friendships", bson.M{}, &friendships, nil)

	if err != nil {
		fmt.Println("error while fetching friendships in bulk: ", err)
	}

	fmt.Println(len(friendships), "friendship seeded!")

	for _, friendship := range friendships {
		for j := 0; j < MESSAGES_COUNTS; j++ {
			from := friendship.UserA

			randNo := rand.IntN(1000)

			if randNo < 500 {
				from = friendship.UserB
			}

			message := mongorm.MessageModel{
				Text:         "Here is the message",
				LastId:       int64(j),
				Sent:         time.Now(),
				FriendshipId: friendship.ID,
				Seen:         time.Now().Add(1000),
				From:         from,
			}

			err := message.Create(context.Background(), db, "messages", &message)

			if err != nil {
				fmt.Println("error while create message: ", message)
			}
		}
	}

	messagesCount, err := mongorm.CountDocuments(context.Background(), db, "messages", bson.M{})

	fmt.Println(messagesCount, "messages seeded!")
}
