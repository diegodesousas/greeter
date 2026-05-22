package database

import (
	"context"

	"github.com/diegodesousas/greeter/internal/domain/greeting"
)

type greetingRepository struct {
	conn Connection
}

func NewGreetingRepository(conn Connection) greeting.Repository {
	return &greetingRepository{conn: conn}
}

func (r *greetingRepository) Save(ctx context.Context, g greeting.Greeting) error {
	_, err := r.conn.Exec(ctx,
		`INSERT INTO greetings (name, greeted_at) VALUES ($1, $2)`,
		g.Name, g.GreetedAt,
	)
	return err
}
