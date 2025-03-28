.PHONY: build
build:
	go build -o ya-boo

.PHONY: docker
docker:
	podman build -t ya-boo:dev .

.PHONY: lint
lint:
	golangci-lint run --config .golangci.yml

.PHONY: statictest
statictest:
	go vet -vettool=$$(which statictest) ./...

.PHONY: up
up:
	podman compose up -d

.PHONY: down
down:
	podman compose down
