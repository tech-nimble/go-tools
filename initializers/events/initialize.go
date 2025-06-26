package events

import "github.com/gookit/event"

func InitializeManager() *event.Manager {
	return event.NewManager("main")
}
