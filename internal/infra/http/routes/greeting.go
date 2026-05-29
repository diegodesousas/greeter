package routes

import (
    "github.com/diegodesousas/go-devkit/pkg/httpserver"
    "github.com/diegodesousas/greeter/internal/application/greet"
    list_greetings "github.com/diegodesousas/greeter/internal/application/list_greetings"
    search_greetings "github.com/diegodesousas/greeter/internal/application/search_greetings"
    "github.com/diegodesousas/greeter/internal/domain/greeting"
    "github.com/diegodesousas/greeter/internal/infra/clock"
    "github.com/diegodesousas/greeter/internal/infra/http/handlers"
)

func Greeting(clock clock.Clock, repo greeting.Repository) []httpserver.Route {
    greetUseCase := greet.NewUseCase(clock, repo)
    listUseCase := list_greetings.NewUseCase(repo)
    searchUseCase := search_greetings.NewUseCase(repo)

    return []httpserver.Route{
        httpserver.NewGet("/hello/{name}", handlers.Hello(greetUseCase)),
        httpserver.NewGet("/greetings", handlers.ListGreetings(listUseCase)),
        httpserver.NewGet("/greetings/search", handlers.SearchGreetings(searchUseCase)),
    }
}
