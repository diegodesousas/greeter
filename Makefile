include .env

NETWORK=diegodesousas-network
PWD=$(shell pwd)
IMAGE_PREFIX=diegodesousas/greeter
MIGRATE_IMAGE=migrate/migrate:v4.18.1
MIGRATE_DB_URL=postgres://${DB_USER}:${DB_PASSWORD}@greeter-postgres:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}
MIGRATE_FLAGS=\
	--rm \
	--name greeter-migrate \
	--network ${NETWORK} \
	-v ${PWD}/migrations:/migrations

DOCKER_DEV_FLAGS = \
	--env-file ${PWD}/.env \
	-v ${PWD}:/application \
	-v greeter-go-mod:/go/pkg/mod \
	-v greeter-go-build-cache:/root/.cache/go-build

.PHONY: build-init build-image build-dev build-network http-up http-down db-up db-down run dev test migrate-up migrate-down migrate-create build-and-run-http

build-init: build-image build-dev build-network

build-image:
	@docker build \
		-t ${IMAGE_PREFIX}-http \
		--target http \
		.

build-dev:
	@docker build \
		-t ${IMAGE_PREFIX}-dev \
		--target dev \
		.

build-network:
ifeq ($(shell docker network ls --format "{{.Name}}" | grep -c "^${NETWORK}$$"),0)
	@echo "Creating network ${NETWORK}"
	@docker network create --driver bridge ${NETWORK}
else
	@echo "Network ${NETWORK} already exists"
endif

http-up:
	@docker run \
		--rm \
		--name greeter-http \
		--network ${NETWORK} \
		--env-file ${PWD}/.env \
		-p ${HTTP_PORT}:${HTTP_PORT} \
		${IMAGE_PREFIX}-http

http-down:
	@echo "Stopping greeter http"
	@docker stop greeter-http > /dev/null
	@echo "HTTP service down"

db-up:
	@docker run \
		--rm \
		-d \
		--name greeter-postgres \
		--network ${NETWORK} \
		-p ${DB_PORT}:5432 \
		-e POSTGRES_USER=${DB_USER} \
		-e POSTGRES_PASSWORD=${DB_PASSWORD} \
		-e POSTGRES_DB=${DB_NAME} \
		postgres:18

db-down:
	@echo "Stopping greeter postgres"
	@docker stop greeter-postgres > /dev/null
	@echo "Postgres down"

dev:
	@docker run \
		--rm \
		--name greeter-dev \
		--network ${NETWORK} \
		-p ${HTTP_PORT}:${HTTP_PORT} \
		${DOCKER_DEV_FLAGS} \
		${IMAGE_PREFIX}-dev \
		go run ./cmd/http/main.go

run:
	@docker run \
		--rm \
		--name greeter-run \
		--network ${NETWORK} \
		-p ${HTTP_PORT}:${HTTP_PORT} \
		${DOCKER_DEV_FLAGS} \
		${IMAGE_PREFIX}-dev \
		go run ./cmd/http/main.go

test:
	@docker run \
		--rm \
		${DOCKER_DEV_FLAGS} \
		${IMAGE_PREFIX}-dev \
		go test ./...

migrate-up:
	@docker run ${MIGRATE_FLAGS} ${MIGRATE_IMAGE} \
		-path=/migrations \
		-database "${MIGRATE_DB_URL}" \
		up

migrate-down:
	@docker run ${MIGRATE_FLAGS} ${MIGRATE_IMAGE} \
		-path=/migrations \
		-database "${MIGRATE_DB_URL}" \
		down 1

migrate-create:
	@docker run ${MIGRATE_FLAGS} \
		--entrypoint migrate ${MIGRATE_IMAGE} \
		create -ext sql -dir /migrations -seq ${NAME}

build-and-run-http: build-image http-up