package http

import (
	"net/http"

	"github.com/pkg/errors"
	"github.com/diegodesousas/go-devkit/pkg/encoding"
	"github.com/diegodesousas/go-devkit/pkg/validator"
)

type ErrorHandlerValidation struct{}

func (e ErrorHandlerValidation) Match(err error) bool {
	_, ok := errors.Cause(err).(validator.Error)

	return ok
}

func (e ErrorHandlerValidation) Write(w http.ResponseWriter, err error) error {
	validatorError := errors.Cause(err).(validator.Error)

	response := DefaultResponse{
		Message: validatorError.Message,
	}

	body, _ := encoding.NewJsonSerializer().Serialize(response)

	w.WriteHeader(http.StatusUnprocessableEntity)

	_, _ = w.Write(body)

	return nil
}

func (e ErrorHandlerValidation) MustLog() bool {
	return false
}
