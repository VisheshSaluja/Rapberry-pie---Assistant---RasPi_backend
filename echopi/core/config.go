package core

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// LoadConfig loads environment variables from .env if available.
func LoadConfig() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("⚠️  Warning: .env not found, using defaults")
	}
}

// GetHotword returns the wake word used for activation (default: "pie")
func GetHotword() string {
	hotword := strings.ToLower(strings.TrimSpace(os.Getenv("HOTWORD")))
	if hotword == "" {
		hotword = "pie"
	}
	return hotword
}

// GetMicDevice returns the FFmpeg input device for recording (default: ":1")
func GetMicDevice() string {
	device := strings.TrimSpace(os.Getenv("INPUT_MIC_DEVICE"))
	if device == "" {
		device = ":1"
	}
	return device
}

// GetModel returns the model used for LLM queries (optional future use)
func GetModel() string {
	model := strings.TrimSpace(os.Getenv("MODEL"))
	if model == "" {
		model = "llama3" // default Ollama model
	}
	return model
}
