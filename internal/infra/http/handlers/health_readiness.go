package handlers

import (
    "fmt"
    stdhttp "net/http"

    infrahttp "github.com/diegodesousas/greeter/internal/infra/http"
    "github.com/diegodesousas/go-devkit/pkg/httpserver"
)

// Pinger checks that a dependency is reachable. database.Connection satisfies it.
type Pinger interface {
    Ping() error
}

// HealthReadiness godoc
//
//	@Summary		Readiness probe
//	@Description	Indica que o serviço está pronto para receber tráfego (banco de dados acessível)
//	@Tags			health
//	@Produce		json
//	@Success		200	{object}	HealthyResponse
//	@Failure		503	{object}	github_com_diegodesousas_greeter_internal_infra_http.DefaultResponse	"Banco de dados indisponível"
//	@Router			/readiness [get]
func HealthReadiness(db Pinger) httpserver.Handler {
    return func(w stdhttp.ResponseWriter, req *stdhttp.Request) error {
        if err := db.Ping(); err != nil {
            return fmt.Errorf("%w: database: %s", infrahttp.ErrServiceUnavailable, err)
        }

        return infrahttp.WriteJson(w, newHealthyResponse())
    }
}
