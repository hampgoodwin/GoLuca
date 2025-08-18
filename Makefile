GOLANGCI_LINT_VERSION = latest

.PHONY: buf
buf: buf-check
	@buf generate proto

bufcheck:
	@buf lint proto
	@buf format proto -w
	@buf breaking proto --against 'https://github.com/hampgoodwin/GoLuca.git#branch=main,ref=HEAD,subdir=proto'


.PHONY: test
test:
	@go test ./... -v --bench . --benchmem --coverprofile=cover.out

testsimple:
	@go test ./... --bench . --benchmem --coverprofile=cover.out

testrace:
	@go test -race ./... -v --bench .

testcovhttp:
	@go test ./... -v --coverprofile=cover.out && go tool cover -html=cover.out

lint:
	@golangci-lint run -v

check: lint test bufcheck vulnerabilitycheck

vulnerabilitycheck:
	@go tool govulncheck ./...

run:
	@go run $$(pwd)/cmd/goluca/main.go

runwiretap:
	@go run $$(pwd)/cmd/wiretap/main.go

up:
	@docker compose -f $$(pwd)/build/package/docker-compose.yml up -d
	@ echo "view jaeger at http://localhost:16686"

down:
	@docker compose -f $$(pwd)/build/package/docker-compose.yml down

downup: down up

dbup:
	@docker compose -f $$(pwd)/build/package/docker-compose.yml up -d db

jaegerup:
	@docker compose -f $$(pwd)/build/package/docker-compose.yml up -d jaeger
	@ echo "view jaeger at http://localhost:16686"

natsup:
	@docker compose -f $$(pwd)/build/package/docker-compose.yml up -d nats
	
# OPEN API COMMANDS
apilint:
	@docker run --rm -v $$PWD/http/v0/spec:/spec redocly/openapi-cli lint /spec/openapi.yml


apipreview:
	redocly preview-docs http/v0/spec/openapi.yml
