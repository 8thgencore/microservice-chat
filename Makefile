MAKEFLAGS += --no-print-directory

# Check if the ENV variable is set
ifneq ($(ENV),)
	include .env.$(ENV)
endif
CONFIG=.env.$(ENV)

# Set the path to the local bin directory
LOCAL_BIN:=$(CURDIR)/bin

# Migration settings
LOCAL_MIGRATION_DIR=$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN="host=localhost port=$(POSTGRES_PORT) dbname=$(POSTGRES_DB) user=$(POSTGRES_USER) password=$(POSTGRES_PASSWORD) sslmode=disable"

# Warning message to ensure correct environment export
.PHONY: check-env
check-env:
ifndef ENV
	$(error "Please run 'export ENV=dev|stage|prod' and 'export $$(xargs < .env.$(ENV))' before executing make")
else 
	@echo "[INFO] Running make with environment: $(ENV)"
endif

# #################### #
# DEPENDENCIES & TOOLS #
# #################### #

_install-global-deps:
	go install github.com/air-verse/air@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest

install-deps: _install-global-deps 
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.35.2
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62.0
	GOBIN=$(LOCAL_BIN) go install mvdan.cc/gofumpt@latest
	GOBIN=$(LOCAL_BIN) go install github.com/gojuno/minimock/v3/cmd/minimock@latest

# Fetch Go dependencies
get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

# Linting
lint:
	GOBIN=$(LOCAL_BIN) bin/golangci-lint run ./... --config .golangci.pipeline.yaml

# Formating
format:
	GOBIN=$(LOCAL_BIN) bin/gofumpt -l -w .

generate-api:
	make generate-chat-api

generate-chat-api:
	mkdir -p pkg/chat/v1
	protoc --proto_path api/chat/v1 --proto_path vendor.protogen \
		--go_out=pkg/chat/v1 --go_opt=paths=source_relative \
		--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go \
		--go-grpc_out=pkg/chat/v1 --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc \
	api/chat/v1/chat.proto

vendor-proto:
		@if [ ! -d vendor.protogen/validate ]; then \
			mkdir -p vendor.protogen/validate &&\
			git clone https://github.com/envoyproxy/protoc-gen-validate vendor.protogen/protoc-gen-validate &&\
			mv vendor.protogen/protoc-gen-validate/validate/*.proto vendor.protogen/validate &&\
			rm -rf vendor.protogen/protoc-gen-validate ;\
		fi

# ##### #
# BUILD #
# ##### #

build-app:
	GOOS=linux GOARCH=amd64 go build -o $(LOCAL_BIN)/main cmd/chat/main.go

docker-net:
	docker network create -d bridge service-net

docker-build: docker-build-app docker-build-migrator

docker-build-app: check-env
	docker buildx build --no-cache --platform linux/amd64 -t chat:${APP_IMAGE_TAG} .

docker-build-migrator: check-env
	docker buildx build --no-cache --platform linux/amd64 -t migrator-chat:${MIGRATOR_IMAGE_TAG} -f migrator.Dockerfile .

# ###### #
# DEPLOY #
# ###### #

docker-deploy: check-env
	docker compose --env-file=.env.$(ENV) up -d

local-migration-status: check-env
	goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} status -v

local-migration-up: check-env
	goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up -v

local-migration-down: check-env
	goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v

# #### #
# STOP #
# #### #

docker-stop: check-env
	docker compose --env-file=.env.$(ENV) down

# ########### #
# DEVELOPMENT #
# ########### #

dev:
	air
