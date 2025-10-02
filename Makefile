ifneq ($(wildcard .env),)
include .env
export
else
$(warning WARNING: .env file not found! Using .env.example)
include .env.example
export
endif

BASE_STACK = docker compose -f docker-compose.yml

# HELP =================================================================================================================
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help

compose-up: ### Run docker compose (without backend and reverse proxy)
	$(BASE_STACK) up --build -d db && docker compose logs -f
.PHONY: compose-up

compose-up-all: ### Run docker compose (with backend and reverse proxy)
	$(BASE_STACK) up --build -d
.PHONY: compose-up-all

compose-down: ### Down docker compose
	$(BASE_STACK) down --remove-orphans
.PHONY: compose-down

swag-v1: ### swag init
	swag init -g internal/controller/http/router.go
.PHONY: swag-v1

proto-v1: ### generate source files from proto
	protoc --go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		docs/proto/v1/*.proto
.PHONY: proto-v1

format: ### Run code formatter
	gofumpt -l -w .
	gci write . --skip-generated -s standard -s default
.PHONY: format

docker-rm-volume: ### remove docker volume
	$(BASE_STACK) down -v --remove-orphans

.PHONY: docker-rm-volume

linter-golangci: ### check by golangci linter
	golangci-lint run
.PHONY: linter-golangci

linter-dotenv: ### check by dotenv linter
	dotenv-linter
.PHONY: linter-dotenv

test: ### run test
	go test -v -race -covermode atomic -coverprofile=coverage.txt ./internal/...
.PHONY: test

run: deps swag-v1 proto-v1 ### swag run for API v1
	go mod download && \
	CGO_ENABLED=0 go run -tags migrate ./cmd/app
.PHONY: run

mock: ### run mockgen
	mockgen -source ./internal/repo/contracts.go -package usecase_test > ./internal/usecase/mocks_repo_test.go
	#mockgen -source ./internal/usecase/contracts.go -package usecase_test > ./internal/usecase/mocks_usecase_test.go
.PHONY: mock

deps: ### deps tidy + verify
	go mod tidy && go mod verify
.PHONY: deps

migrate-create:  ### create new migration
	migrate create -ext sql -dir migrations ${NAME}
.PHONY: migrate-create

migrate-up: ### migration up
	migrate -path migrations -database '$(PG_URL)?sslmode=disable' up
.PHONY: migrate-up

bin-deps: ### install tools
	go install tool
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate
.PHONY: bin-deps

pre-commit: swag-v1 proto-v1 mock format linter-golangci test ### run pre-commit
.PHONY: pre-commit


#psql "host=localhost port=5432 dbname=db user=user password='myAwEsOm3pa55@w0rd' sslmode=disable"

