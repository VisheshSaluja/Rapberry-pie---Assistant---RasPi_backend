package plugins

import (
	"fmt"
	"strings"

	"echopi/core/common"
)

type SystemSkill struct{}

func init() { Register(&SystemSkill{}) }

func (s *SystemSkill) CanHandle(input string) bool {
	return strings.Contains(input, "sleep") ||
		strings.Contains(input, "stop listening") ||
		strings.Contains(input, "goodbye")
}

func (s *SystemSkill) Handle(input string) error {
	common.Speak("Okay, going to sleep now.")
	fmt.Println("😴 Entering sleep mode...")
	return fmt.Errorf("sleep")
}
