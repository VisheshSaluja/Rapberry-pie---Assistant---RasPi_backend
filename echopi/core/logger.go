package core

import (
	"fmt"
	"os"
	"time"
)

func LogInteraction(user, assistant string) error {
	logDir := "data"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		os.MkdirAll(logDir, 0755)
	}

	f, err := os.OpenFile("data/conversations.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("log open error: %v", err)
	}
	defer f.Close()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	entry := fmt.Sprintf("[%s]\nUSER: %s\nECHO: %s\n\n",
		timestamp, user, assistant)
	_, err = f.WriteString(entry)
	return err
}
