package routes

import (
	stdhttp "net/http"

	"github.com/diegodesousas/go-devkit/pkg/httpserver"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func Docs() []httpserver.Route {
	uiHandler := httpSwagger.Handler(httpSwagger.URL("/docs/doc.json"))

	adapter := func(w stdhttp.ResponseWriter, req *stdhttp.Request) error {
		uiHandler.ServeHTTP(w, req)
		return nil
	}

	return []httpserver.Route{
		httpserver.NewGet("/docs/*", adapter),
	}
}
