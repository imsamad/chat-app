package main

import (
	"chat-app/handlers"
	"chat-app/middlewares"
	"chat-app/mongorm"
	"chat-app/utils"
	"encoding/json"

	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

type Handler func(w http.ResponseWriter, r *http.Request) *utils.Response

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if response := h(w, r); response != nil {
		w.Header().Set("Content-Type", "application/json")
		if response.Code != 0 {
			w.WriteHeader(response.Code)
		}
		if response.Message != "" {
			json.NewEncoder(w).Encode(response.Message)
		}
	}
}

var (
	MONGODB_URI string
)

func InitEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("no .env file found")
	}
	MONGODB_URI = os.Getenv("MONGODB_URI")

	if MONGODB_URI == "" {
		fmt.Println("MONGODB_URI is reqiored")
		os.Exit(0)
	}

}

func init() {
	InitEnv()
}

func runServer() {
	mongoClient, err := mongorm.Connect(MONGODB_URI)

	if err != nil {
		fmt.Println("error while connecting to mongo: \n reason ", err)
		os.Exit(1)
	}

	fmt.Println("DB connected!")

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	r.Use(middleware.Logger)

	db := mongoClient.Database("chat-db")
	handler := handlers.New(db)

	authRouter := chi.NewRouter()
	authRouter.Post("/signup", Handler(handler.SignupController).ServeHTTP)
	authRouter.Post("/login", Handler(handler.LoginHandler).ServeHTTP)
	authRouter.Post("/logout", Handler(handler.Logout).ServeHTTP)

	authRouter.With(middlewares.AuthMiddleware(db)).Get("/me", Handler(handler.GetMe).ServeHTTP)

	friendsRouter := chi.NewRouter()
	friendsRouter.Use(middlewares.AuthMiddleware(db))
	friendsRouter.Post("/", Handler(handler.AddFriend).ServeHTTP)
	friendsRouter.Delete("/", Handler(handler.RemoveFriend).ServeHTTP)
	friendsRouter.Get("/", Handler(handler.GetMyFriends).ServeHTTP)

	messagesRouters := chi.NewRouter()
	messagesRouters.Use(middlewares.AuthMiddleware(db))
	messagesRouters.Post("/", Handler(handler.CreateMessage).ServeHTTP)

	r.With(middlewares.AuthMiddleware(db)).Get("/init-state", Handler(handler.GetInitialState).ServeHTTP)
	r.Mount("/auth", authRouter)
	r.Mount("/friends", friendsRouter)
	r.Mount("/messages", messagesRouters)

	err = http.ListenAndServe(":4000", r)

	if err != nil {
		fmt.Println(err)
		fmt.Println("unable to run application")
		os.Exit(1)
	}
}
