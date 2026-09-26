package greet

import (
	"context"
	"fmt"
	"unicode/utf8"

	"github.com/diegodesousas/go-devkit/pkg/validator"
)

func nameRequired(_ context.Context, dto DTO) error {
	if validator.IsEmpty(dto.Name) {
		return validator.NewRequiredError("name")
	}

	return nil
}

func nameMaxLength(_ context.Context, dto DTO) error {
	var maxLength = 50

	if utf8.RuneCountInString(dto.Name) > maxLength {
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
