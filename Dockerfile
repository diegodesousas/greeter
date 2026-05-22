FROM golang:1.26-alpine AS base
RUN apk update && apk add --no-cache git pkgconf gcc libc-dev
WORKDIR /application
COPY go.mod go.sum ./
RUN go mod download

FROM base AS dev

FROM base AS builder
COPY . .
RUN CGO_ENABLED=1 go build -tags musl -o /application/bin/app /application/cmd/http/main.go

FROM alpine:3.23 AS http
WORKDIR /application
EXPOSE 3000
COPY --from=builder /application/bin/app ./bin/app
CMD ["/application/bin/app"]
