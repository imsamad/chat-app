package handlers

import (
	"chat-app/mongorm"
	"chat-app/utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type Claims struct {
	jwt.RegisteredClaims
	ID string `json:"id,omitempty"`
}

func (h *Handler) SignupController(w http.ResponseWriter, r *http.Request) *utils.Response {

	var body User
	json.NewDecoder(r.Body).Decode(&body)

	message := map[string]string{
		"email":    "",
		"password": "",
		"name":     "",
	}

	if body.Email == "" {
		message["email"] = "Email is required"
	}
	if body.Password == "" {
		message["password"] = "Password is required"
	}
	if body.Name == "" {
		message["name"] = "Name is required"
	}

	if message["email"] != "" || message["password"] != "" || message["name"] != "" {
		return &utils.Response{
			Code:    http.StatusBadRequest,
			Message: message,
		}
	}

	var user mongorm.UserModel

	user.Read(context.Background(), h.DB, "users", bson.M{"email": body.Email}, &user)

	if user.Email != "" {
		return &utils.Response{
			Code: http.StatusConflict,
			Message: map[string]string{
				"email": fmt.Sprintf("user with this email already exist"),
			},
		}
	}
	user.Email = body.Email
	user.Name = body.Name
	user.Friends = []primitive.ObjectID{}
	hashedPwd, err := utils.HashPassword(body.Password)

	if err != nil {
		fmt.Println("below is error during hashing the passwrd")
		fmt.Println(err)
		return &utils.Response{
			Message: map[string]string{"message": "Please try again"},
			Code:    http.StatusInternalServerError,
		}
	}

	user.Password = hashedPwd

	user.Create(context.Background(), h.DB, "users", &user)
	oneMonth := time.Now().AddDate(0, 1, 0)
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(oneMonth),
			Subject:   user.ID.Hex(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		fmt.Println("error occuring during jwt-fying token ==>", err)
		return &utils.Response{
			Message: "Please try again",
			Code:    http.StatusInternalServerError,
		}
	}
	w.Header().Set("Content-Type", "application/json")

	secure := true
	if os.Getenv("APP_ENV") == "development" {
		secure = false
	}
	cookie := &http.Cookie{
		Name:     "user",
		Value:    tokenStr,
		Expires:  oneMonth,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, cookie)

	return &utils.Response{
		Code: 200,
		Message: map[string]interface{}{
			// "jwt": tokenStr,
			"user": map[string]string{
				"email": user.Email,
				"name":  user.Name,
			},
		},
	}
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) *utils.Response {

	var body User

	json.NewDecoder(r.Body).Decode(&body)

	if body.Email == "" || body.Password == "" {
		return &utils.Response{
			Code: http.StatusBadRequest,
			Message: map[string]string{
				"email":    "Email is required",
				"password": "Password is required",
			},
		}
	}

	var user mongorm.UserModel

	err := user.Read(context.Background(), h.DB, "users", bson.M{"email": body.Email}, &user)

	if err != nil {
		fmt.Println(err)
		return &utils.Response{
			Code: http.StatusConflict,
			Message: map[string]string{
				"email": "Email not found",
			},
		}
	}

	if utils.CheckPasswordHash(user.Password, body.Password) != true {
		return &utils.Response{
			Code: http.StatusConflict,
			Message: map[string]string{
				"password": "Password is incorrect",
			},
		}
	}

	oneMonth := time.Now().AddDate(0, 1, 0)
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(oneMonth),
			Subject:   user.ID.Hex(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		fmt.Println("error occuring during jwt-fying token ==>", err)
		return &utils.Response{
			Message: "Please try again",
			Code:    http.StatusInternalServerError,
		}
	}

	secure := true
	if os.Getenv("APP_ENV") == "development" {
		secure = false
	}

	cookie := &http.Cookie{
		Name:     "user",
		Value:    tokenStr,
		Expires:  oneMonth,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, cookie)
	return &utils.Response{
		Code: 200,
		Message: map[string]interface{}{
			// "jwt": tokenStr,
			"user": map[string]string{
				"email": user.Email,
				"name":  user.Name,
			},
		},
	}

}

func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) *utils.Response {
	user_id := r.Context().Value("user").(primitive.ObjectID)

	var user mongorm.UserModel

	err := user.Read(context.Background(), h.DB, "users", bson.M{"_id": user_id}, &user)

	if err != nil {
		return &utils.Response{
			Code:    http.StatusNoContent,
			Message: "There is not user",
		}
	}

	return &utils.Response{
		Code: http.StatusOK,
		Message: struct {
			ID    string `json:"id"`
			Email string `json:"email"`
			Name  string `json:"name"`
		}{
			ID:    user.ID.Hex(),
			Email: user.Email,
			Name:  user.Name,
		},
	}
}
