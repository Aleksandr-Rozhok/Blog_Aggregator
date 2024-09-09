package models

import (
	"blog_aggregator/internal/database"
	"github.com/google/uuid"
	"time"
)

type Status struct {
	Status string `json:"status"`
}

type ApiConfig struct {
	DB *database.Queries
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	APIKey    string    `json:"api_key"`
}
