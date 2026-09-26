package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/diegodesousas/greeter/internal/infra/sanitize"
)

// ReadJson decodes the request body into dest and trims surrounding whitespace
// from every string field (opt out with the `trim:"-"` tag). An empty body is
// not an error, leaving validation of required fields to the use case.
// Malformed JSON is wrapped with ErrBadRequest so the error handler maps it to 400.
func ReadJson(req *http.Request, dest any) error {
	if err := json.NewDecoder(req.Body).Decode(dest); err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("%w: %s", ErrBadRequest, err)
	}

	sanitize.TrimStrings(dest)

	return nil
}
