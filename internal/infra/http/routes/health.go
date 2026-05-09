package routes

import (
    "github.com/diegodesousas/go-devkit/pkg/httpserver"
    "github.com/diegodesousas/greeter/internal/infra/http/handlers"
)

func Health() []httpserver.Route {
    return []httpserver.Route{
        httpserver.NewGet("/liveness", handlers.HealthLiveness()),
        httpserver.NewGet("/readiness", handlers.HealthReadiness()),
    }
}
