package tools

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	w.WriteHeader(code)

	dat, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		fmt.Printf("Error writing response: %v\n", err)
	}

	write, err := w.Write(dat)
	if err != nil {
		fmt.Printf("Error writing respones: %s", err)
	}

	fmt.Printf("Response written to: %d bytes\n", write)
}
