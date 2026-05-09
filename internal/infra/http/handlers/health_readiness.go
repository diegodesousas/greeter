package handlers

import (
    stdhttp "net/http"

    infrahttp "github.com/diegodesousas/greeter/internal/infra/http"
    "github.com/diegodesousas/go-devkit/pkg/httpserver"
)

func HealthReadiness() httpserver.Handler {
    return func(w stdhttp.ResponseWriter, req *stdhttp.Request) error {
        return infrahttp.WriteJson(w, newHealthyResponse())
    }
}
