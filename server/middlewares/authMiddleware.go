package middlewares

import (
	"chat-app/mongorm"
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func verifyToken(tokenString string) (jwt.MapClaims, error) {
	secretKey := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func AuthMiddleware(db *mongo.Database) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth_token := r.Header.Get("Authorization")

			if auth_token == "" {
				w.WriteHeader(404)
				w.Write([]byte("You are not authorised"))
				return
			}

			token, err := verifyToken(auth_token)

			if err != nil {
				w.WriteHeader(404)
				w.Write([]byte("You are not authorised"))
				return
			}

			user_id_hex, err := token.GetSubject()
			if err != nil {
				w.WriteHeader(404)
				w.Write([]byte("You are not authorised"))
				return
			}

			user_id, err := primitive.ObjectIDFromHex(user_id_hex)

			if err != nil {
				w.WriteHeader(404)
				w.Write([]byte("You are not authorised"))
				return
			}

			var user mongorm.UserModel

			err = user.Read(context.Background(), db, "users", bson.M{"_id": user_id}, &user)

			if err != nil {
				fmt.Println("error while fetching user from db: ", err)
				w.WriteHeader(404)
				w.Write([]byte("You are not authorised"))
				return
			}

			ctx := context.WithValue(r.Context(), "user", user_id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
