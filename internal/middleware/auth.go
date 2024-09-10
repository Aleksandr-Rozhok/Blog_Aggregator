package middleware

import (
	"blog_aggregator/internal/database"
	"blog_aggregator/models"
	"blog_aggregator/tools"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"strings"
	"time"
)

type ApiConfig struct {
	DB *database.Queries
}

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *ApiConfig) MiddlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		apiKey := r.Header.Get("Authorization")
		if strings.HasPrefix(apiKey, "ApiKey ") {
			apiKey = strings.TrimPrefix(apiKey, "ApiKey ")
		} else {
			http.Error(w, "Missing or invalid API key", http.StatusUnauthorized)
			return
		}

		user, err := cfg.DB.GetUserByAPIKey(ctx, apiKey)
		if err != nil {
			fmt.Printf("Error getting user: %v\n", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		handler(w, r, user)
	}
}

func (cfg *ApiConfig) HandlerUsersGet(w http.ResponseWriter, r *http.Request, user database.User) {
	tools.RespondWithJSON(w, http.StatusOK, user)
}

func (cfg *ApiConfig) HandlerFeedsPost(w http.ResponseWriter, r *http.Request, user database.User) {
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

	var feedsData models.Feed
	if err := json.Unmarshal(bodyBytes, &feedsData); err != nil {
		fmt.Printf("Error unmarshalling body: %v\n", err)
		return
	}

	if feedsData.Name == "" {
		fmt.Errorf("Name can't be blank")
		return
	}

	feeds, err := cfg.DB.CreateFeed(ctx, database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedsData.Name,
		Url:       feedsData.URL,
		UserID:    user.ID,
	})
	if err != nil {
		fmt.Printf("Error creating user: %v\n", err)
	}

	tools.RespondWithJSON(w, http.StatusOK, feeds)
}
