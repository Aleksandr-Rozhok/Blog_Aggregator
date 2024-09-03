package tools

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Error struct {
	Error string `json:"error"`
}

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	w.WriteHeader(code)

	resBody := Error{
		Error: msg,
	}

	dat, err := json.MarshalIndent(resBody, "", "  ")
	if err != nil {
		fmt.Printf("Error writing response: %v\n", err)
	}

	write, err := w.Write(dat)
	if err != nil {
		fmt.Printf("Error writing respones: %s", err)
	}

	fmt.Printf("Response written to: %d bytes\n", write)
}
