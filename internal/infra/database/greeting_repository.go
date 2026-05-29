package database

import (
	"context"
	"fmt"
	"time"

	"github.com/diegodesousas/greeter/internal/domain/greeting"
)

type greetingModel struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	GreetedAt time.Time `db:"greeted_at"`
}

func (m greetingModel) toDomain() greeting.Greeting {
	return greeting.Greeting{
		ID:        m.ID,
		Name:      m.Name,
		Message:   fmt.Sprintf("Hello, %s!", m.Name),
		GreetedAt: m.GreetedAt,
	}
}

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

func (r *greetingRepository) Search(ctx context.Context, name string, page, perPage int) ([]greeting.Greeting, int, error) {
	pattern := "%" + name + "%"

	var total int
	if err := r.conn.Get(ctx, &total,
		`SELECT COUNT(*) FROM greetings WHERE unaccent(name) ILIKE unaccent($1)`,
		pattern,
	); err != nil {
		return nil, 0, err
	}

	var models []greetingModel
	offset := (page - 1) * perPage
	if err := r.conn.Select(ctx, &models,
		`SELECT id, name, greeted_at FROM greetings WHERE unaccent(name) ILIKE unaccent($1) ORDER BY greeted_at DESC LIMIT $2 OFFSET $3`,
		pattern, perPage, offset,
	); err != nil {
		return nil, 0, err
	}

	greetings := make([]greeting.Greeting, len(models))
	for i, m := range models {
		greetings[i] = m.toDomain()
	}

	return greetings, total, nil
}

func (r *greetingRepository) List(ctx context.Context, page, perPage int) ([]greeting.Greeting, int, error) {
	var total int
	if err := r.conn.Get(ctx, &total, `SELECT COUNT(*) FROM greetings`); err != nil {
		return nil, 0, err
	}

	var models []greetingModel
	offset := (page - 1) * perPage
	if err := r.conn.Select(ctx, &models,
		`SELECT id, name, greeted_at FROM greetings ORDER BY greeted_at DESC LIMIT $1 OFFSET $2`,
		perPage, offset,
	); err != nil {
		return nil, 0, err
	}

	greetings := make([]greeting.Greeting, len(models))
	for i, m := range models {
		greetings[i] = m.toDomain()
	}

	return greetings, total, nil
}
