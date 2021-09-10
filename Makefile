test:
	go test ./... -v --bench . --benchmem --covermode=count

lint:
	docker run --rm -v $$(pwd):/app -w /app golangci/golangci-lint:v1.42.1 golangci-lint run -v

check: lint test

up:
	docker-compose -f $$(pwd)/build/package/docker-compose.yml up -d

down:
	docker-compose -f $$(pwd)/build/package/docker-compose.yml down

dbup:
	docker-compose -f $$(pwd)/build/package/docker-compose.yml up -d db