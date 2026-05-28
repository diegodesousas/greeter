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

// HealthLiveness godoc
//
//	@Summary		Liveness probe
//	@Description	Indica que o processo está vivo (não verifica dependências externas)
//	@Tags			health
//	@Produce		json
//	@Success		200	{object}	HealthyResponse
//	@Router			/liveness [get]
func HealthLiveness() httpserver.Handler {
    return func(w stdhttp.ResponseWriter, req *stdhttp.Request) error {
        return infrahttp.WriteJson(w, newHealthyResponse())
    }
}
