package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"echopi/core/plugins"
)

func main() {
	fmt.Println("🧠 Starting EchoPi backend server on Pi...")

	// Force plugin init (ensures all init() ran)
	_ = plugins.List()

	// Test endpoint
	http.HandleFunc("/api/test", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"message": "EchoPi backend running",
		})
	})

	// Main skill router
	http.HandleFunc("/api/skill", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "POST only", http.StatusMethodNotAllowed)
			return
		}

		input := r.URL.Query().Get("q")
		if input == "" {
			http.Error(w, "missing q parameter", http.StatusBadRequest)
			return
		}

		// Use unified plugin routing
		handled, result := plugins.TryHandleUnified(input)

		if !handled {
			json.NewEncoder(w).Encode(map[string]string{
				"status": "unhandled",
			})
			return
		}

		// Send structured JSON response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})

	log.Println("🧠 EchoPi backend listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
