# Backlog da revisão geral do greeter

Especificação para retomar a revisão geral do projeto feita em 2026-09-26. Os itens 1 a 7 foram corrigidos; este documento descreve **o que falta**, com evidência no código, proposta e critério de aceite de cada item, além do processo e das decisões já tomadas.

A numeração dos itens segue a revisão original e é usada nas descrições dos PRs ("Item N da revisão geral do projeto").

## Concluído

| Item | O quê | PR |
|---|---|---|
| 1 | JSON inválido em `POST /hello` retorna 400 (`ErrBadRequest`, `infrahttp.ReadJson`) | #11 |
| 2 | Defaults de `page`/`per_page` e 422 para valor não numérico (`infrahttp.QueryInt`, `DefaultPage`/`DefaultPerPage`) | #13 |
| 3 | Limite do nome conta caracteres, não bytes | #14 |
| 4 | Trim de toda string da requisição na borda HTTP (`sanitize.TrimStrings`, `QueryString`, `PathParam`) | #16 (substituiu #15) |
| 5 | Escape de `%`/`_`/`\` na busca e limite de 50 no `name` da busca (`validation.ExceedsMaxLength`) | #17 |
| 6 | Readiness checa o banco (503 via `ErrServiceUnavailable`) | #18 |
| 7 | Shutdown com prazo e fechamento do banco (`shutdown.Graceful` do go-devkit) | #20 |
| — | Mocks gerados com mockery no pacote `mocks/` | #19 |
| — | go-devkit v0.0.1 → v0.3.0 e Go 1.27 | #21 |
| — | `pkg/shutdown` no go-devkit, release v0.3.0 | go-devkit#26 |

## Processo combinado

- **Um item por PR**, branch a partir da `main` atualizada (`fix/`, `feature/`, `chore/`), PR contra `main`. O merge é feito só quando pedido, com **squash**.
- **TDD**: teste falhando pelo motivo certo antes do código de produção.
- **Tudo em Docker**: `make test`, `make mocks`, `make docs`; `go build`/`go vet` via `docker run ... diegodesousas/greeter-dev`.
- **Mocks**: nunca escritos à mão. Adicionar a interface em `.mockery.yml` e rodar `make mocks`.
- **Verificação de ponta a ponta** quando o comportamento muda: subir a branch num container temporário na porta 3001 (`--network diegodesousas-network --env-file .env -e HTTP_PORT=3001`), sem mexer no servidor da porta 3000. Para cenários destrutivos (banco caindo, etc.), usar um Postgres descartável (`postgres:18` com usuário `verify`), nunca o banco de dev.
- **Mudou uma convenção?** Atualizar o `CLAUDE.md` no mesmo PR.

### Armadilhas conhecidas do ambiente

- O shell é **zsh**: `$var` sem aspas **não** é dividido em palavras, `path` é uma variável especial (ligada ao `PATH`) e `--include=*.go` sem aspas falha por glob. Use `--include='*.go'` e nomes de variáveis como `u`/`p`.
- `gh pr edit` falha (erro de "Projects (classic)"). Para editar a descrição de um PR: `gh api -X PATCH repos/diegodesousas/greeter/pulls/N -F "body=@arquivo.md"`.
- Merge via API, fixando o SHA conferido: `gh api -X PUT repos/diegodesousas/greeter/pulls/N/merge -f merge_method=squash -f sha=<sha>`.
- Force push é bloqueado pelo ambiente. Para desfazer algo já enviado, usar `git revert`.
- O `mergeable_state` do GitHub às vezes fica `unknown`. Conferir localmente com `git merge-base --is-ancestor origin/main HEAD`.

## Decisões já tomadas (não rediscutir)

- **Trim de strings acontece na borda HTTP**, não nos use cases (#16). Handlers leem entrada só pelos helpers do `infrahttp`.
- **Mocks centralizados** em `mocks/`, gerados pelo mockery v3 (template testify, `unroll-variadic: true`).
- **Código genérico e reutilizável entre serviços vai para o go-devkit** (ex.: `pkg/shutdown`), que tem convenções próprias (`github.com/pkg/errors`, `assert`, mocks à mão, `make release`).
- **Refactors pequenos e literais**: quando foi pedido para generalizar a checagem de tamanho, a versão com regra genérica (`validation.MaxLength[T]`) e constante de domínio foi **desfeita**. A escolhida foi o helper simples `validation.ExceedsMaxLength(value, max) bool`, com a constante declarada dentro da função que a usa. Siga esse estilo.

## Pendentes

A ordem abaixo é a sugerida: primeiro os rápidos e os que destravam os seguintes (CI), depois arquitetura, depois performance e operação.

### Item 8 — `.env.example` e URL de migração quebram o ambiente Docker

- **Evidência:** `.env.example:11` tem `DB_HOST=localhost`, mas a aplicação roda num container na rede `diegodesousas-network`, onde o banco é `greeter-postgres`. `Makefile:7` monta `MIGRATE_DB_URL` com `greeter-postgres:${DB_PORT}`: `DB_PORT` é a porta publicada no host, e dentro da rede a porta é sempre `5432`.
- **Proposta:** `DB_HOST=greeter-postgres` no `.env.example`; `MIGRATE_DB_URL` com a porta fixa `5432`.
- **Aceite:** com um `.env` copiado do exemplo, `make db-up`, `make migrate-up` e `make dev` funcionam, inclusive com `DB_PORT` diferente de 5432.

### Item 21 — Documentação e arquivos desatualizados

- **Evidência:**
  - `README.md:9` descreve `GET /hello?name=<name>` e fala em um único endpoint.
  - `insomnia.yaml:10` chama `GET /hello/thays`. Uma correção existe na branch `chore/insomnia-post-hello` (PR #12, fechado sem merge a pedido).
  - `.idea/` está versionado, apesar de estar no `.gitignore`.
  - `make dev` e `make run` são idênticos.
  - `include .env` (`Makefile:1`) quebra o Makefile se o `.env` ainda não existe.
  - O agente se chama `.claude/agents/senior-enginner.md` (typo).
- **Proposta:** README com os endpoints atuais (`POST /hello`, `GET /greetings`, `GET /greetings/search`, health, `/docs`) e os alvos do Makefile; `git rm --cached -r .idea`; unificar `dev`/`run`; `-include .env`; renomear o agente. **Confirmar com o usuário** antes de reaproveitar ou descartar a branch do Insomnia, porque o PR #12 foi fechado de propósito.
- **Aceite:** README reflete a API real; `git ls-files .idea` vazio; `make` sem `.env` mostra erro claro em vez de falhar no include.

### Item 20 — Sem CI e sem lint

- **Evidência:** não existe `.github/`. `gofmt -l cmd internal` lista 19 arquivos (indentados com espaços).
- **Proposta:** workflow do GitHub Actions com `gofmt -l` (falhando se houver saída), `go vet ./...`, `golangci-lint` e `go test ./...`, em Go 1.27. Formatar os 19 arquivos num PR próprio antes de ligar a checagem.
- **Aceite:** PRs rodam CI automaticamente; `gofmt -l` vazio.
- **Nota:** o go-module `go-devkit` é público, então o CI não precisa de credenciais para baixá-lo.

### Item 19 — Arquivos sem `_test.go` correspondente

- **Evidência:** sem teste: `internal/infra/database/repositories.go`, `internal/infra/http/pagination.go` (só constantes), `internal/infra/http/routes/docs.go`, `routes/greeting.go`, `routes/health.go`. `domain/greeting/repository.go` é só interface.
- **Proposta:** testes das rotas (método, path e handler registrados; `/docs/*` servindo a UI) e de `NewRepositories`. Interfaces e arquivos só de constantes podem ficar sem teste. Registrar essa exceção no `CLAUDE.md`.
- **Aceite:** toda lógica tem teste; a exceção fica documentada.

### Item 9 — Application depende de infra

- **Evidência:** `internal/application/greet/greet.go:8` importa `internal/infra/clock`, violando a regra de dependência da Clean Architecture.
- **Proposta:** a interface `Clock` passa a viver em application (ou domain); `internal/infra/clock` mantém só a implementação. Ajustar `.mockery.yml` para o novo pacote e regenerar os mocks.
- **Aceite:** `grep -rn "internal/infra" internal/application internal/domain` vazio.

### Item 10 — Regra de negócio duplicada

- **Evidência:** `fmt.Sprintf("Hello, %s!", ...)` em `internal/domain/greeting/greeting.go:18` e em `internal/infra/database/greeting_repository.go:22` (`toDomain`).
- **Proposta:** o repositório reconstrói a entidade por uma função do domínio (ex.: `greeting.Restore(id, name, greetedAt)` ou reaproveitando `greeting.New` e atribuindo o ID), de modo que a mensagem tenha uma única fonte.
- **Aceite:** a string `Hello, %s!` aparece uma única vez no código de produção.

### Item 12 — 404 conhece `sql.ErrNoRows`

- **Evidência:** `internal/infra/http/error_handler_not_found.go:17` casa `sql.ErrNoRows`, que é detalhe do banco.
- **Proposta:** erro de domínio `greeting.ErrNotFound` (ou um `ErrNotFound` genérico), traduzido no repositório; o writer de 404 casa o erro de domínio.
- **Aceite:** `internal/infra/http` não importa `database/sql`.
- **Nota:** hoje nenhum fluxo retorna "não encontrado". Implementar junto com o primeiro endpoint que precise, ou só a troca de dependência.

### Item 13 — Inconsistências menores da API

- **Evidência:**
  - `internal/infra/http/validation_error_handler.go:26` cria um serializer a cada chamada e ignora o erro (`body, _ :=`), enquanto os outros writers usam `WriteJson`.
  - A resposta 422 devolve só `message` e perde o `code` do `validator.Error`.
  - `POST /hello` responde 200, e não 201, sem o `id`, embora o `GreetingDTO` da listagem tenha `id`.
- **Proposta:** o writer de validação usa `WriteJson` e inclui `code`; `POST /hello` responde 201 com `id`. Isso exige que o `Save` devolva o ID, por exemplo com `INSERT ... RETURNING id`.
- **Aceite:** 422 com `{"message": ..., "code": ...}`; `POST /hello` → 201 com `id`, Swagger regenerado.
- **Atenção:** muda o contrato da API (status e corpo). Confirmar com o usuário.

### Item 11 — `list_greetings` e `search_greetings` quase idênticos

- **Evidência:** `GreetingDTO`, `PaginationDTO` e `Output` duplicados em `internal/application/list_greetings/list_greetings.go:16-31` e `search_greetings/search_greetings.go:17-32`. As regras de paginação (`page >= 1`, `1 <= per_page <= 100`) e o mapeamento também se repetem.
- **Proposta (escolher com o usuário):** (a) um único `GET /greetings?name=` com filtro opcional, removendo `/greetings/search`, o que muda a API; ou (b) extrair um pacote comum de paginação e DTOs, sem mudar a API.
- **Aceite:** sem DTOs e regras de paginação duplicados.

### Item 22 — Nomes de pacote com underscore

- **Evidência:** `internal/application/list_greetings`, `search_greetings`. Não é idiomático em Go, e os imports precisam de alias.
- **Proposta:** `listgreetings`, `searchgreetings`. Fazer depois do item 11, que pode eliminar um deles, e atualizar `.mockery.yml`.
- **Aceite:** nenhum pacote com underscore.

### Item 14 — Busca sempre faz varredura completa

- **Evidência:** as migrations não criam índices. `unaccent(name) ILIKE '%x%'` não usa índice, e não há índice em `greeted_at` para o `ORDER BY`.
- **Proposta:** migration com uma função `immutable_unaccent` (o `unaccent` não é `IMMUTABLE`), `pg_trgm` e índice GIN sobre `immutable_unaccent(name)`; índice em `greeted_at DESC`. A query de busca passa a usar a mesma função.
- **Aceite:** `EXPLAIN` da busca e da listagem usa os índices; testes de integração (item 18) cobrindo a query.

### Item 18 — Queries nunca testadas contra um Postgres real

- **Evidência:** `internal/infra/database/greeting_repository_test.go` usa `MockConnection`; `unaccent`, `ILIKE`, escape e paginação só foram conferidos manualmente.
- **Proposta:** testes de integração com build tag (padrão do go-devkit: `//go:build integration`), subindo Postgres (dockertest ou container irmão), rodando as migrations, e um alvo `make test-integration`.
- **Aceite:** `make test-integration` cobre `Save`, `List` e `Search`, incluindo acentos e curingas.

### Item 16 — Exposição pública sem proteção

- **Evidência:** `cmd/http/main.go:74` usa `httpserver.AllowAll()` (CORS aberto); `POST /hello` é público, sem rate limit; o body não tem limite de tamanho (sem `http.MaxBytesReader`); o Swagger (`routes.Docs()`, `cmd/http/main.go:102`) é exposto em qualquer ambiente.
- **Proposta:** limitar o body no `ReadJson`, com `MaxBytesReader` e 413 ou 400; CORS configurável por ambiente (o go-devkit só tem `AllowAll`, então pode exigir uma option nova lá); Swagger só fora de produção (`ENV`); rate limit por IP, avaliando se ele vai para o go-devkit.
- **Aceite:** body acima do limite é rejeitado; `/docs` indisponível com `ENV=production`.

### Item 17 — Dockerfile de produção

- **Evidência:** a imagem `http` (`Dockerfile:15`) roda como root; `CGO_ENABLED=1 -tags musl` (`Dockerfile:13`) sem necessidade aparente, porque o pgx é Go puro, embora o dd-trace possa ter motivo para isso; as migrations não fazem parte do deploy.
- **Proposta:** `USER` não root na imagem final; investigar se o CGO é mesmo necessário e, se não for, removê-lo; definir como as migrations rodam no deploy (job separado com `migrate/migrate` ou etapa de CI/CD).
- **Aceite:** a imagem roda como usuário não root e o smoke test da imagem de produção continua passando.

### Item 15 — `COUNT(*)` + `OFFSET`

- **Evidência:** `List` e `Search` fazem duas queries e paginam com `OFFSET`. Funciona em volume pequeno e escala mal.
- **Proposta:** só quando houver volume: paginação por cursor (keyset em `greeted_at`, `id`). Muda a API.
- **Aceite:** decisão registrada; implementar só se necessário.

## Pendências que surgiram durante a sessão

- **`nameRequired` duplicado** em `greet_validator.go:11` e `search_greetings_validator.go:11`. Candidato a helper no mesmo estilo de `validation.ExceedsMaxLength`.
- **Constante local em `greet_validator.go:20`:** `var maxLength = 50`. A busca já usa `const` dentro da função; alinhar.
- **Readiness sem timeout próprio:** `Ping()` do go-devkit não recebe `context` (`health_readiness.go:27`). Melhoria no go-devkit: `PingContext(ctx)`.
- **Branches antigas não removidas** (local e remoto), já mergeadas ou fechadas: `chore/dev-run-command`, `chore/insomnia-post-hello`, `chore/mockery-centralized-mocks`, `chore/upgrade-go-devkit-v0.3.0`, `feature/list-greetings-paginated`, `feature/postgres-config`, `feature/refactor-hello-to-post`, `feature/swagger-api-docs`, `feature/trim-request-strings`, `fix/escape-search-like-wildcards`, `fix/graceful-shutdown-timeout`, `fix/invalid-json-body-returns-400`, `fix/name-max-length-counts-characters`, `fix/pagination-query-defaults`, `fix/postgres-data-persistence`, `fix/readiness-checks-database`, `fix/trim-greet-name`. Confirmar antes de apagar `chore/insomnia-post-hello` (ver item 21).
