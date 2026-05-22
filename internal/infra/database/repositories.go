package database

import "github.com/diegodesousas/greeter/internal/domain/greeting"

type Repositories struct {
	Greeting greeting.Repository
}

func NewRepositories(conn Connection) Repositories {
	return Repositories{
		Greeting: NewGreetingRepository(conn),
	}
}
