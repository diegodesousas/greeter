package search_greetings

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
	const maxNameLength = 50

	if validation.ExceedsMaxLength(dto.Name, maxNameLength) {
		return validator.Error{
			Code:    "max_length",
			Message: fmt.Sprintf("attribute name must have at most %d characters", maxNameLength),
		}
	}
	return nil
}

func pageMin1(_ context.Context, dto DTO) error {
	if dto.Page < 1 {
		return validator.Error{
			Code:    "min_value",
			Message: "attribute page must be at least 1",
		}
	}
	return nil
}

func perPageBetween1And100(_ context.Context, dto DTO) error {
	if dto.PerPage < 1 || dto.PerPage > 100 {
		return validator.Error{
			Code:    "out_of_range",
			Message: fmt.Sprintf("attribute per_page must be between 1 and 100"),
		}
	}
	return nil
}

func newSearchGreetingsValidator() validator.Validator[DTO] {
	return validator.New[DTO](
		nameRequired,
		nameMaxLength,
		pageMin1,
		perPageBetween1And100,
	)
}
