.PHONY: build test lint docker run-local migrate clean

build:
	go build -o bin/ra ./cmd/ra

test:
	go test ./... -v

lint:
	golangci-lint run ./...

docker:
	docker build -t gigvault/ra:local .

run-local: docker
	../infra/scripts/deploy-local.sh ra

migrate:
	migrate -path migrations -database $$DATABASE_URL up

clean:
	rm -rf bin/
	go clean

