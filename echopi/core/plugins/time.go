package plugins

import (
	"fmt"
	"strings"
	"time"

	"echopi/core/common"
)

type TimeSkill struct{}

func init() { Register(&TimeSkill{}) }

func (t *TimeSkill) CanHandle(input string) bool {
	return strings.Contains(input, "time") || strings.Contains(input, "date")
}

func (t *TimeSkill) Handle(input string) error {
	if strings.Contains(input, "time") {
		msg := fmt.Sprintf("It's %s.", time.Now().Format("3:04 PM"))
		fmt.Println("⏰ " + msg)
		common.Speak(msg)
		return nil
	}
	if strings.Contains(input, "date") {
		msg := fmt.Sprintf("Today is %s.", time.Now().Format("Monday, January 2, 2006"))
		fmt.Println("📆 " + msg)
		common.Speak(msg)
		return nil
	}
	return nil
}
