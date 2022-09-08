-include .makerc
export

run: main.go
	go run main.go --host 0.0.0.0 --port 8080

install: go.mod
	go install gitlab.com/so_literate/genmock/cmd/genmock@latest
	go mod tidy

start: build up

build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down --volumes -t 0

lint:
	go vet ./...

mocks:
	go generate ./...

test:
	go test ./...

test-coverage:
	go test -covermode=count -coverpkg=./... -coverprofile coverage.out -v ./...
	go tool cover -html coverage.out -o coverage.html
