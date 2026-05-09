package greeting

import (
	"fmt"
	"time"
)

type Greeting struct {
	Name      string
	Message   string
	GreetedAt time.Time
}

func New(name string, now time.Time) Greeting {
	return Greeting{
		Name:      name,
		Message:   fmt.Sprintf("Hello, %s!", name),
		GreetedAt: now,
	}
}
