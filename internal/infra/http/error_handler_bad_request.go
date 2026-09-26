package http

import (
	"errors"
	"net/http"
)

var ErrBadRequest = errors.New("bad request")

type ErrorHandlerBadRequest struct{}

func (e ErrorHandlerBadRequest) MustLog() bool {
	return false
}

func (e ErrorHandlerBadRequest) Match(err error) bool {
	return errors.Is(err, ErrBadRequest)
}

func (e ErrorHandlerBadRequest) Write(w http.ResponseWriter, _ error) error {
	w.WriteHeader(http.StatusBadRequest)

	response := DefaultResponse{
		Message: "invalid request body",
	}

	return WriteJson(w, response)
}
