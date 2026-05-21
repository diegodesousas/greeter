include .env

NETWORK=pppay-network
PWD=$(shell pwd)

build-init: build-image build-dev build-network

build-image:
	@docker build \
		-t pppay/greeter-http \
		--target http \
		--no-cache \
		.

build-dev:
	@docker build \
		-t pppay/greeter-dev \
		--target base \
		.

build-network:
ifeq ($(shell docker network list --filter name=${NETWORK} | wc -l), 1)
	@echo "Creating network ${NETWORK}"
	@docker network create --driver bridge ${NETWORK}
else
	@echo "Network ${NETWORK} created"
endif

http-up:
	@docker run \
		--rm \
		--name greeter-http \
		--network ${NETWORK} \
		--env-file ${PWD}/.env \
		-p ${HTTP_PORT}:${HTTP_PORT} \
		pppay/greeter-http

http-down:
	@echo "Stopping greeter http"
	@docker stop greeter-http > /dev/null
	@echo "HTTP service down"

run:
	@docker run \
		--rm \
		--name greeter-run \
		--network ${NETWORK} \
		--env-file ${PWD}/.env \
		-p ${HTTP_PORT}:${HTTP_PORT} \
		-v ${PWD}:/application \
		-v greeter-go-mod:/go/pkg/mod \
		-v greeter-go-build-cache:/root/.cache/go-build \
		pppay/greeter-dev \
		go run ./cmd/http/main.go

test:
	@docker run \
		--rm \
		--env-file ${PWD}/.env \
		-v ${PWD}:/application \
		-v greeter-go-mod:/go/pkg/mod \
		-v greeter-go-build-cache:/root/.cache/go-build \
		pppay/greeter-dev \
		go test ./...

build-and-run-http: build-image http-up
