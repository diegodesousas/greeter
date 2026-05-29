package search_greetings

import (
	"context"
	"time"

	"github.com/diegodesousas/go-devkit/pkg/validator"
	"github.com/diegodesousas/greeter/internal/domain/greeting"
)

type DTO struct {
	Name    string
	Page    int
	PerPage int
}

type GreetingDTO struct {
	ID        string    `json:"id"`
	Message   string    `json:"message"`
	GreetedAt time.Time `json:"greeted_at"`
}

type PaginationDTO struct {
	Total   int `json:"total"`
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

type Output struct {
	Data       []GreetingDTO `json:"data"`
	Pagination PaginationDTO `json:"pagination"`
}

type UseCase interface {
	Run(ctx context.Context, dto DTO) (Output, error)
}

type searchGreetingsUseCase struct {
	repo      greeting.Repository
	validator validator.Validator[DTO]
}

func NewUseCase(repo greeting.Repository) UseCase {
	return searchGreetingsUseCase{
		repo:      repo,
		validator: newSearchGreetingsValidator(),
	}
}

func (u searchGreetingsUseCase) Run(ctx context.Context, dto DTO) (Output, error) {
	if err := u.validator.Validate(ctx, dto); err != nil {
		return Output{}, err
	}

	greetings, total, err := u.repo.Search(ctx, dto.Name, dto.Page, dto.PerPage)
	if err != nil {
		return Output{}, err
	}

	data := make([]GreetingDTO, len(greetings))
	for i, g := range greetings {
		data[i] = GreetingDTO{
			ID:        g.ID,
			Message:   g.Message,
			GreetedAt: g.GreetedAt,
		}
	}

	return Output{
		Data: data,
		Pagination: PaginationDTO{
			Total:   total,
			Page:    dto.Page,
			PerPage: dto.PerPage,
		},
	}, nil
}
