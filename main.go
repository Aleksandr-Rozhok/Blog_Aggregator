package main

import (
	"blog_aggregator/models"
	"blog_aggregator/tools"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mux := http.NewServeMux()
	port := os.Getenv("PORT")

	server := &http.Server{
		Addr:    port,
		Handler: mux,
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

	err = http.ListenAndServe(server.Addr, server.Handler)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Server is running")
	}
}
