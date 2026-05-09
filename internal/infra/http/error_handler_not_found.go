package http

import (
	"database/sql"
	"net/http"

	"github.com/pkg/errors"
)

type ErrorHandlerNotFound struct{}

func (e ErrorHandlerNotFound) MustLog() bool {
	return false
}

func (e ErrorHandlerNotFound) Match(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

func (e ErrorHandlerNotFound) Write(w http.ResponseWriter, _ error) error {
	w.WriteHeader(http.StatusNotFound)

	response := DefaultResponse{
		Message: "resource not found",
	}

	return WriteJson(w, response)
}
