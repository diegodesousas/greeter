package handlers

import (
    stdhttp "net/http"

    infrahttp "github.com/diegodesousas/greeter/internal/infra/http"
    "github.com/diegodesousas/go-devkit/pkg/httpserver"
)

type HealthyResponse struct {
    Message string `json:"message"`
}

func newHealthyResponse() HealthyResponse {
    return HealthyResponse{Message: "ok"}
}

func HealthLiveness() httpserver.Handler {
    return func(w stdhttp.ResponseWriter, req *stdhttp.Request) error {
        return infrahttp.WriteJson(w, newHealthyResponse())
    }
}
