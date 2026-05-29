package greeting

import "context"

type Repository interface {
	Save(ctx context.Context, g Greeting) error
	List(ctx context.Context, page, perPage int) ([]Greeting, int, error)
	Search(ctx context.Context, name string, page, perPage int) ([]Greeting, int, error)
}
