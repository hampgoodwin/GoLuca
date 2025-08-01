GOLANGCI_LINT_VERSION = latest

.PHONY: buf
buf:
	@ sh ./scripts/buf.sh
	@buf lint proto
	@buf format proto -w
	@buf generate proto

.PHONY: test
test:
	go test ./... -v --bench . --benchmem --coverprofile=cover.out

testsimple:
	go test ./... --bench . --benchmem --coverprofile=cover.out

testrace:
	go test -race ./... -v --bench .

testcovhttp:
	go test ./... -v --coverprofile=cover.out && go tool cover -html=cover.out

lint:
	@docker run --rm -t -v $(shell pwd):/app -w /app \
	--user $(shell id -u):$(shell id -g) \
	-v $(shell go env GOCACHE):/.cache/go-build -e GOCACHE=/.cache/go-build \
	-v $(shell go env GOMODCACHE):/.cache/mod -e GOMODCACHE=/.cache/mod \
	-v $(HOME)/.cache/golangci-lint:/.cache/golangci-lint -e GOLANGCI_LINT_CACHE=/.cache/golangci-lint \
	golangci/golangci-lint:$(GOLANGCI_LINT_VERSION) golangci-lint run -v

check: lint test

run:
	go run $$(pwd)/cmd/goluca/main.go
runwiretap:
	go run $$(pwd)/cmd/wiretap/main.go

up:
	docker compose -f $$(pwd)/build/package/docker-compose.yml up -d
	@ echo "view jaeger at http://localhost:16686"

down:
	docker compose -f $$(pwd)/build/package/docker-compose.yml down

downup: down up

dbup:
	docker compose -f $$(pwd)/build/package/docker-compose.yml up -d db

jaegerup:
	docker compose -f $$(pwd)/build/package/docker-compose.yml up -d jaeger
	@ echo "view jaeger at http://localhost:16686"

natsup:
	docker compose -f $$(pwd)/build/package/docker-compose.yml up -d nats
	
# OPEN API COMMANDS
apilint:
	docker run --rm -v $$PWD/http/v0/spec:/spec redocly/openapi-cli lint /spec/openapi.yml


apipreview:
	sh ./scripts/openapi_previewdocs.sh
