package http

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/diegodesousas/go-devkit/pkg/httpserver"
	"github.com/diegodesousas/go-devkit/pkg/validator"
)

// QueryString reads a string query param with surrounding whitespace removed.
// A missing param returns an empty string.
func QueryString(req *http.Request, key string) string {
	return strings.TrimSpace(req.URL.Query().Get(key))
}

// PathParam reads a route path param with surrounding whitespace removed.
// A missing param returns an empty string.
func PathParam(req *http.Request, key string) string {
	return strings.TrimSpace(httpserver.GetParam(req, key))
}

// QueryInt reads an integer query param. A missing or empty param returns
// fallback; a non-integer value returns a validator.Error so the error handler
// maps it to 422. Range checks are left to the use case validators.
func QueryInt(req *http.Request, key string, fallback int) (int, error) {
	raw := req.URL.Query().Get(key)
	if raw == "" {
		return fallback, nil
	}

	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, validator.Error{
			Code:    "invalid_type",
			Message: fmt.Sprintf("attribute %s must be an integer", key),
		}
	}

	return value, nil
}
