# Greeter

## Responsibility

Example service demonstrating the Clean Architecture structure used across diegodesousas-com-br services. Exposes a single greeting endpoint.

## Features

- `GET /hello?name=<name>` — returns a greeting message

## Develop environment

- Copy `.env.example` to `.env` and fill in your credentials
- To run locally: `make run`
- To run with Docker: `make build-init` then `make http-up`
