FROM golang:1.24-alpine AS base
COPY . /application
RUN apk update && apk add --no-cache git pkgconf gcc libc-dev
WORKDIR /application

FROM base AS builder
RUN CGO_ENABLED=1 go build -tags musl -o /application/bin/app /application/cmd/http/main.go

FROM alpine:3.18 AS http
EXPOSE 3000
COPY --from=builder /application/bin/app /application/bin/app
CMD ["./application/bin/app"]
