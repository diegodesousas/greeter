package greeting

import "context"

type Repository interface {
	Save(ctx context.Context, g Greeting) error
}
