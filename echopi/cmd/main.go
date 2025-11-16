package main

import (
	"fmt"
	"log"
	"os"

	"echopi/core"

	"github.com/joho/godotenv"
)

func main() {
	// 1️⃣ Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		log.Println("⚠️  Warning: .env file not found, using defaults")
	}

	// 2️⃣ Check for continuous mode
	if len(os.Args) > 1 && os.Args[1] == "loop" {
		fmt.Println("🔊 Initializing EchoPi Phase 7: Continuous Listening Mode...")
		core.StartContinuousListening()
		fmt.Println("✅ EchoPi Phase 7 complete.")
		return
	}

	// 3️⃣ Default: one-shot pipeline
	fmt.Println("🔊 Initializing EchoPi Phase 4 (single run)...")
	if err := core.RunEchoPi(); err != nil {
		log.Fatalf("❌ EchoPi pipeline failed: %v", err)
	}

	fmt.Println("✅ EchoPi Phase 4 complete.")
}
