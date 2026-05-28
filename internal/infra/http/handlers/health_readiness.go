package handlers

import (
    stdhttp "net/http"

    infrahttp "github.com/diegodesousas/greeter/internal/infra/http"
    "github.com/diegodesousas/go-devkit/pkg/httpserver"
)

// HealthReadiness godoc
//
//	@Summary		Readiness probe
//	@Description	Indica que o serviço está pronto para receber tráfego
//	@Tags			health
//	@Produce		json
//	@Success		200	{object}	HealthyResponse
//	@Router			/readiness [get]
func HealthReadiness() httpserver.Handler {
    return func(w stdhttp.ResponseWriter, req *stdhttp.Request) error {
        return infrahttp.WriteJson(w, newHealthyResponse())
    }
}
