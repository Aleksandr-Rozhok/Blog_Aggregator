package main

import (
	"blog_aggregator/internal/database"
	"blog_aggregator/internal/middleware"
	"blog_aggregator/models"
	"blog_aggregator/tools"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mux := http.NewServeMux()
	port := os.Getenv("PORT")
	dbURL := os.Getenv("CONNECTION_STRING")
	db, err := sql.Open("postgres", dbURL)
	dbQueries := database.New(db)

	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	config := middleware.ApiConfig{
		DB: dbQueries,
	}

	mux.HandleFunc("GET /v1/healthz", func(w http.ResponseWriter, r *http.Request) {
		statusOK := models.Status{
			Status: "OK",
		}

		tools.RespondWithJSON(w, http.StatusOK, statusOK)
	})
	mux.HandleFunc("GET /v1/err", func(w http.ResponseWriter, r *http.Request) {
		tools.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
	})
	mux.HandleFunc("POST /v1/users", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				fmt.Printf("Error writing response: %v\n", err)
			}
		}(r.Body)

		var userData models.User
		if err := json.Unmarshal(bodyBytes, &userData); err != nil {
			fmt.Printf("Error unmarshalling body: %v\n", err)
			return
		}

		if userData.Name == "" {
			fmt.Errorf("Name can't be blank")
			return
		}

		user, err := config.DB.CreateUser(ctx, database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      userData.Name,
		})
		if err != nil {
			fmt.Printf("Error creating user: %v\n", err)
		}

		tools.RespondWithJSON(w, http.StatusOK, user)
	})
	mux.Handle("GET /v1/users", config.MiddlewareAuth(config.HandlerUsersGet))
	mux.Handle("POST /v1/feeds", config.MiddlewareAuth(config.HandlerFeedsPost))
	mux.HandleFunc("GET /v1/feeds", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		feeds, err := config.DB.GetFeeds(ctx)
		if err != nil {
			fmt.Printf("Error getting feeds: %v\n", err)
		}

		fmt.Println(feeds)

		tools.RespondWithJSON(w, http.StatusOK, feeds)
	})

	err = http.ListenAndServe(server.Addr, server.Handler)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Server is running")
	}
}
