ifneq (,$(wildcard ./.env))
    include .env
    export
endif

APP?=${APP_NAME}
REGISTRY?=rg.nl-ams.scw.cloud/main
COMMIT_SHA=$(shell git rev-parse --short HEAD)

.PHONY: build
## build: build the application (default for make file)
build: clean
	@echo " -- Building"
	@go build -o ./build/${APP} -v ./cmd/bot

.DEFAULT_GOAL := build

.PHONY: run
## run: runs go run main.go
run:
	go run -race ./cmd/bot/main.go

.PHONY: clean
## clean: cleans the binary
clean:
	@echo " -- Cleaning"
	@go clean

.PHONY: test
## test: runs go test with default values
test:
	go test -v -count=1 -race ./...

.PHONY: docker-build
## docker-build: builds the stringifier docker image to registry
docker-build: build
	docker build -t ${APP}:${COMMIT_SHA} .

.PHONY: docker-push
## docker-push: pushes the stringifier docker image to registry
docker-push: check-environment docker-build
	@echo " -- In error when push check docker login rg.nl-ams.scw.cloud/main -u nologin -p $SCW_SECRET_TOKEN"
	docker push ${REGISTRY}/${ENV}/${APP}:${COMMIT_SHA}

.PHONY: help
## help: Prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

# helper rule for deployment
check-environment:
	@echo "Checking ${ENV} file "
