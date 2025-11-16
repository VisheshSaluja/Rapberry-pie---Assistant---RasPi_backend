package core

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"echopi/core/plugins"
)

type BackendRequest struct {
	Input string `json:"input"`
}

type BackendResponse struct {
	Output string `json:"output"`
}

// StartBackendServer starts a simple HTTP server on the Pi
func StartBackendServer() {
	port := os.Getenv("ECHOPI_BACKEND_PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/process", handleProcess)

	addr := ":" + port
	log.Printf("🧠 EchoPi backend listening on %s\n", addr)

	s := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	if err := s.ListenAndServe(); err != nil {
		log.Fatalf("❌ backend server failed: %v", err)
	}
}

func handleProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}

	var req BackendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	input := req.Input
	if input == "" {
		json.NewEncoder(w).Encode(BackendResponse{Output: "I didn’t receive any text."})
		return
	}

	// 1️⃣ Try plugins first (notes, alarms, time, etc.)
	handled, err := plugins.TryHandle(input)
	if err != nil {
		out := fmt.Sprintf("Plugin error: %v", err)
		_ = json.NewEncoder(w).Encode(BackendResponse{Output: out})
		return
	}
	if handled {
		// Many plugins speak directly via TTS; you can also make them return strings instead.
		_ = json.NewEncoder(w).Encode(BackendResponse{Output: ""})
		return
	}

	// 2️⃣ Fallback to LLM (optional: you can keep LLM local on Mac if you prefer)
	resp, err := QueryLLM(input)
	if err != nil {
		_ = json.NewEncoder(w).Encode(BackendResponse{
			Output: "I had trouble thinking about that right now.",
		})
		return
	}

	_ = json.NewEncoder(w).Encode(BackendResponse{Output: resp})
}
