package greet

import (
    "context"
    "fmt"

    "github.com/diegodesousas/go-devkit/pkg/validator"
)

func nameRequired(_ context.Context, dto DTO) error {
    if validator.IsEmpty(dto.Name) {
        return validator.NewRequiredError("name")
    }

    return nil
}

func nameMaxLength(_ context.Context, dto DTO) error {
    if len(dto.Name) > 50 {
        return validator.Error{
            Code:    "max_length",
            Message: fmt.Sprintf("attribute name must have at most 50 characters"),
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
