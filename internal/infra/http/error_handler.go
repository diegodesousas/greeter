package http

import (
	"context"
	"net/http"

	"github.com/diegodesousas/go-devkit/pkg/httpserver"
	"github.com/diegodesousas/go-devkit/pkg/log"
)

type ErrorWriter interface {
	Match(err error) bool
	Write(w http.ResponseWriter, err error) error
	MustLog() bool
}

type DefaultResponse struct {
	Message string `json:"message"`
}

func defaultWrite(ctx context.Context, w http.ResponseWriter, originalErr error) {
	w.WriteHeader(http.StatusInternalServerError)

	response := DefaultResponse{
		Message: "internal server error",
	}

	body, _ := serializer.Serialize(response)

	_, _ = w.Write(body)

	log.Error(ctx, originalErr)
}

func NewErrorHandler(list ...ErrorWriter) httpserver.ErrorHandler {
	return func(ctx context.Context, w http.ResponseWriter, err error) {
		for _, writer := range list {
			if writer.Match(err) {
				if err := writer.Write(w, err); err != nil {
					defaultWrite(ctx, w, err)
					return
				}

				if !writer.MustLog() {
					return
				}

				log.Error(ctx, err)
				return
			}
		}

		defaultWrite(ctx, w, err)
	}
}
