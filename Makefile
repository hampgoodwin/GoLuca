.PHONY: test
test:
	go test ./... -v --bench . --benchmem --coverprofile=count

testrace:
	go test ./... -v --bench . --benchmen -race

testcovhttp:
	go test ./... --coverprofile=cover.out && go tool cover -html=cover.out

lint:
	docker run --rm -v $$(pwd):/app -w /app golangci/golangci-lint:v1.42.1 golangci-lint run -v

check: lint test

run:
	go run $$(pwd)/cmd/goluca/main.go

up:
	docker-compose -f $$(pwd)/build/package/docker-compose.yml up -d

down:
	docker-compose -f $$(pwd)/build/package/docker-compose.yml down

downup: down up

dbup:
	docker-compose -f $$(pwd)/build/package/docker-compose.yml up -d db

# OPEN API COMMANDS
apilint:
	docker run --rm -v $$PWD/api:/spec redocly/openapi-cli lint /spec/openapi.yml

apipreview:
	sh ./scripts/openapi_previewdocs.sh