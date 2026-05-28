package list_greetings

import (
	"context"
	"fmt"

	"github.com/diegodesousas/go-devkit/pkg/validator"
)

func pageRequired(_ context.Context, dto DTO) error {
	if dto.Page < 1 {
		return validator.Error{
			Code:    "min_value",
			Message: "attribute page must be at least 1",
		}
	}
	return nil
}

func perPageRange(_ context.Context, dto DTO) error {
	if dto.PerPage < 1 || dto.PerPage > 100 {
		return validator.Error{
			Code:    "out_of_range",
			Message: fmt.Sprintf("attribute per_page must be between 1 and 100"),
		}
	}
	return nil
}

func newListGreetingsValidator() validator.Validator[DTO] {
	return validator.New[DTO](
		pageRequired,
		perPageRange,
	)
}
