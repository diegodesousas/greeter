package routes

import (
    "github.com/diegodesousas/go-devkit/pkg/httpserver"
    "github.com/diegodesousas/greeter/internal/application/greet"
    "github.com/diegodesousas/greeter/internal/infra/clock"
    "github.com/diegodesousas/greeter/internal/infra/http/handlers"
)

func Greeting(clock clock.Clock) []httpserver.Route {
    useCase := greet.NewUseCase(clock)

    return []httpserver.Route{
        httpserver.NewGet("/hello/{name}", handlers.Hello(useCase)),
    }
}
