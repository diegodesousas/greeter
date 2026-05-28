package handlers

import (
    "net/http"
    "strconv"

    "github.com/diegodesousas/go-devkit/pkg/httpserver"
    "github.com/diegodesousas/greeter/internal/application/greet"
    list_greetings "github.com/diegodesousas/greeter/internal/application/list_greetings"
    infrahttp "github.com/diegodesousas/greeter/internal/infra/http"
)

// Hello godoc
//
//	@Summary		Cumprimenta uma pessoa pelo nome
//	@Description	Gera uma saudação para o nome informado e persiste o registro no histórico
//	@Tags			greetings
//	@Produce		json
//	@Param			name	path		string	true	"Nome da pessoa (1-50 caracteres)"
//	@Success		200		{object}	greet.GreetingDTO
//	@Failure		422		{object}	github_com_diegodesousas_greeter_internal_infra_http.DefaultResponse	"Erro de validação"
//	@Failure		500		{object}	github_com_diegodesousas_greeter_internal_infra_http.DefaultResponse	"Erro interno"
//	@Router			/hello/{name} [get]
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

// ListGreetings godoc
//
//	@Summary		Lista saudações paginadas
//	@Description	Retorna a lista de saudações registradas, com paginação
//	@Tags			greetings
//	@Produce		json
//	@Param			page		query		int	false	"Número da página (>= 1)"	default(1)
//	@Param			per_page	query		int	false	"Itens por página (1-100)"	default(10)
//	@Success		200			{object}	list_greetings.Output
//	@Failure		422			{object}	github_com_diegodesousas_greeter_internal_infra_http.DefaultResponse	"Erro de validação"
//	@Failure		500			{object}	github_com_diegodesousas_greeter_internal_infra_http.DefaultResponse	"Erro interno"
//	@Router			/greetings [get]
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
