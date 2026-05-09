package greet

import (
    "context"
    "time"

    "github.com/diegodesousas/greeter/internal/domain/greeting"
    "github.com/diegodesousas/greeter/internal/infra/clock"
    "github.com/diegodesousas/go-devkit/pkg/validator"
)

type DTO struct {
    Name string
}

type GreetingDTO struct {
    Message   string    `json:"message"`
    GreetedAt time.Time `json:"greeted_at"`
}

type UseCase interface {
    Run(ctx context.Context, dto DTO) (GreetingDTO, error)
}

type greetUseCase struct {
    clock     clock.Clock
    validator validator.Validator[DTO]
}

func NewUseCase(clock clock.Clock) UseCase {
    return greetUseCase{
        clock:     clock,
        validator: newGreetValidator(),
    }
}

func (u greetUseCase) Run(ctx context.Context, dto DTO) (GreetingDTO, error) {
    if err := u.validator.Validate(ctx, dto); err != nil {
        return GreetingDTO{}, err
    }

    g := greeting.New(dto.Name, u.clock.Now())

    return GreetingDTO{
        Message:   g.Message,
        GreetedAt: g.GreetedAt,
    }, nil
}
