package greeting

import (
	"fmt"
	"time"
)

// NameMaxLength is the maximum number of characters a greeted name may have.
const NameMaxLength = 50

type Greeting struct {
	ID        string
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
