package handlers

import (
    "net/http"

    "github.com/diegodesousas/go-devkit/pkg/httpserver"
    "github.com/diegodesousas/greeter/internal/application/greet"
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
