package core

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

// RecordAudio records `duration` seconds of microphone input into `outputFile`
func RecordAudio(outputFile string, duration int) error {
	// Use INPUT_MIC_DEVICE from .env, fallback to default Mac mic (:1)
	device := os.Getenv("INPUT_MIC_DEVICE")
	if device == "" {
		device = ":1"
	}

	fmt.Printf("🎙️  Recording %ds of audio from device %s...\n", duration, device)

	cmd := exec.Command(
		"ffmpeg",
		"-f", "avfoundation",
		"-i", device,
		"-ar", "16000", // resample to 16kHz for Whisper
		"-ac", "1",     // mono
		"-t", fmt.Sprintf("%d", duration),
		"-filter:a", "volume=2.0", // amplify gain
		"-y", outputFile, // overwrite existing
	)

	// Show FFmpeg logs if DEBUG=true
	if os.Getenv("DEBUG") == "true" {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	start := time.Now()
	if err := cmd.Run(); err != nil {
		// Graceful handling for Ctrl+C or mic interruptions
		if _, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("❌ failed to record audio: %v", err)
		}
	}
	fmt.Printf("✅  Audio captured successfully (%.3fs)\n", time.Since(start).Seconds())
	return nil
}
