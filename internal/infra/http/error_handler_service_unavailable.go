package http

import (
	"errors"
	"net/http"
)

var ErrServiceUnavailable = errors.New("service unavailable")

type ErrorHandlerServiceUnavailable struct{}

// MustLog is true so an unavailable dependency (e.g. database down) shows up
// in the logs, not only as a failing probe.
func (e ErrorHandlerServiceUnavailable) MustLog() bool {
	return true
}

func (e ErrorHandlerServiceUnavailable) Match(err error) bool {
	return errors.Is(err, ErrServiceUnavailable)
}

func (e ErrorHandlerServiceUnavailable) Write(w http.ResponseWriter, _ error) error {
	w.WriteHeader(http.StatusServiceUnavailable)

	response := DefaultResponse{
		Message: "service unavailable",
	}

	return WriteJson(w, response)
}
