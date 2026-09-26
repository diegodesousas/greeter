package greet

import (
	"context"
	"fmt"

	"github.com/diegodesousas/go-devkit/pkg/validator"
	"github.com/diegodesousas/greeter/internal/application/validation"
)

func nameRequired(_ context.Context, dto DTO) error {
	if validator.IsEmpty(dto.Name) {
		return validator.NewRequiredError("name")
	}

	return nil
}

func nameMaxLength(_ context.Context, dto DTO) error {
	var maxLength = 50

	if validation.ExceedsMaxLength(dto.Name, maxLength) {
		return validator.Error{
			Code:    "max_length",
			Message: fmt.Sprintf("attribute name must have at most %d characters", maxLength),
		}
	}

	return nil
}

func newGreetValidator() validator.Validator[DTO] {
	return validator.New[DTO](
		nameRequired,
		nameMaxLength,
	)
}
