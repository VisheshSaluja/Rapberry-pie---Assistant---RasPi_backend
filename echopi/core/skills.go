//package core

//import "echopi/core/plugins"

//func HandleSkill(transcript string) (bool, error) {
//    return plugins.TryHandle(transcript)
//}






package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"echopi/core/plugins"
)

// ---------------------------------------------------------------------
// UI EVENT CHANNEL (goroutine-safe)
// ---------------------------------------------------------------------

var uiChannel = make(chan map[string]interface{}, 20)

// Called by plugins to emit UI updates (e.g., Notes, Pomodoro countdown)
func EmitUI(payload map[string]interface{}) {
	uiChannel <- payload
}

// The backend loop listens to UIChannel() and returns UI messages upward
func UIChannel() chan map[string]interface{} {
	return uiChannel
}

// ---------------------------------------------------------------------
// HandleSkill: Execute plugins AND return UI + Speak via panic→parse
// ---------------------------------------------------------------------

func HandleSkill(transcript string) (bool, string, map[string]interface{}, error) {
	defer func() {
		// Recover from plugin panics — used to return UI or speak events
		if r := recover(); r != nil {
			fmt.Println("🔥 SkillEngine panic intercepted:", r)
		}
	}()

	handled, err := plugins.TryHandle(transcript)
	if err != nil {
		// The plugin returns UI and Speak through specially formatted error:
		//     "ui::{json}|speak::{msg}"
		return handled, parseSpeak(err), parseUI(err), nil
	}

	return handled, "", nil, nil
}

// ---------------------------------------------------------------------
// Parse "speak::..." out of the error
// ---------------------------------------------------------------------

func parseSpeak(err error) string {
	if err == nil {
		return ""
	}
	parts := strings.Split(err.Error(), "|")
	for _, p := range parts {
		if strings.HasPrefix(p, "speak::") {
			return strings.TrimPrefix(p, "speak::")
		}
	}
	return ""
}

// ---------------------------------------------------------------------
// Parse "ui::{json}" out of the error
// ---------------------------------------------------------------------

func parseUI(err error) map[string]interface{} {
	if err == nil {
		return nil
	}
	parts := strings.Split(err.Error(), "|")
	for _, p := range parts {
		if strings.HasPrefix(p, "ui::") {
			raw := strings.TrimPrefix(p, "ui::")

			var payload map[string]interface{}
			if json.Unmarshal([]byte(raw), &payload) == nil {
				return payload
			}
		}
	}
	return nil
}

// ---------------------------------------------------------------------
// Listen for UI events produced by goroutines (Pomodoro countdown)
// This is called by main loop in backend API handler.
// ---------------------------------------------------------------------

func DrainUIEvents() (map[string]interface{}, error) {
	select {
	case ui := <-uiChannel:
		return ui, nil
	default:
		return nil, errors.New("no_ui_events")
	}
}
