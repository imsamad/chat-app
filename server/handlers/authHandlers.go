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

	// jsonify user creds
	var body User
	json.NewDecoder(r.Body).Decode(&body)

	message := map[string]string{
		"email":    "",
		"password": "",
		"name":     "",
	}

	// validate data
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

	// check if user already exist with the user sent email
	err := user.Read(context.Background(), h.DB, "users", bson.M{"email": body.Email}, user, nil)

	if err != nil || user.Email != "" {
		return &utils.Response{
			Code: http.StatusConflict,
			Message: map[string]string{
				"email": fmt.Sprintf("user with this email already exist"),
			},
		}
	}
	user.Email = body.Email
	user.Name = body.Name

	hashedPwd, err := utils.HashPassword(body.Password)

	if err != nil {
		return &utils.Response{
			Message: map[string]string{"message": "Please try again"},
			Code:    http.StatusInternalServerError,
		}
	}

	// set hashed pwd
	user.Password = hashedPwd

	// create user entry
	err = user.Create(context.Background(), h.DB, "users", &user)
	if err != nil {
		return &utils.Response{
			Message: map[string]string{"message": "Please try again"},
			Code:    http.StatusInternalServerError,
		}
	}

	// issue jwt token
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
			"jwt": tokenStr,
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

	err := user.Read(context.Background(), h.DB, "users", bson.M{"email": body.Email}, &user, nil)

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
			"jwt": tokenStr,
			"user": map[string]string{
				"email": user.Email,
				"name":  user.Name,
			},
		},
	}
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) *utils.Response {
	http.SetCookie(w, &http.Cookie{
		Name:   "user",
		Value:  "",
		MaxAge: -1,
	})

	return &utils.Response{
		Code: 200,
		Message: map[string]interface{}{
			"message": "logged out successfully!",
		},
	}
}

func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) *utils.Response {
	user_id := r.Context().Value("user").(primitive.ObjectID)
	fmt.Print("i was hitted")
	var user mongorm.UserModel

	err := user.Read(context.Background(), h.DB, "users", bson.M{"_id": user_id}, &user, nil)

	if err != nil {
		fmt.Println("err: ", err)
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
