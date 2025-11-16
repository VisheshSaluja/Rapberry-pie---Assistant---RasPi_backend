package main

import (
	"log"

	"echopi/core"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	text := "EchoPi is ready to speak!"
	if err := core.Speak(text); err != nil {
		log.Fatal(err)
	}
}
