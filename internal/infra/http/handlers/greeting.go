package handlers

import (
    "net/http"
    "strconv"

    "github.com/diegodesousas/go-devkit/pkg/httpserver"
    "github.com/diegodesousas/greeter/internal/application/greet"
    list_greetings "github.com/diegodesousas/greeter/internal/application/list_greetings"
    infrahttp "github.com/diegodesousas/greeter/internal/infra/http"
)

func Hello(useCase greet.UseCase) httpserver.Handler {
    return func(w http.ResponseWriter, req *http.Request) error {
        dto := greet.DTO{
            Name: httpserver.GetParam(req, "name"),
        }

        result, err := useCase.Run(req.Context(), dto)
        if err != nil {
            return err
        }

        return infrahttp.WriteJson(w, result)
    }
}

func ListGreetings(useCase list_greetings.UseCase) httpserver.Handler {
    return func(w http.ResponseWriter, req *http.Request) error {
        page, _ := strconv.Atoi(req.URL.Query().Get("page"))
        perPage, _ := strconv.Atoi(req.URL.Query().Get("per_page"))

        dto := list_greetings.DTO{
            Page:    page,
            PerPage: perPage,
        }

        result, err := useCase.Run(req.Context(), dto)
        if err != nil {
            return err
        }

        return infrahttp.WriteJson(w, result)
    }
}
