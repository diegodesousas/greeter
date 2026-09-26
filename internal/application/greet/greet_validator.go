package greet

import (
	"context"

	"github.com/diegodesousas/go-devkit/pkg/validator"
	"github.com/diegodesousas/greeter/internal/application/validation"
	"github.com/diegodesousas/greeter/internal/domain/greeting"
)

func nameRequired(_ context.Context, dto DTO) error {
	if validator.IsEmpty(dto.Name) {
		return validator.NewRequiredError("name")
	}

	return nil
}

func newGreetValidator() validator.Validator[DTO] {
	return validator.New[DTO](
		nameRequired,
		validation.MaxLength("name", greeting.NameMaxLength, func(dto DTO) string { return dto.Name }),
	)
}
