install:
	go mod download
	go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	go install github.com/vektra/mockery/v2@v2.46.1

gen:
	go generate ./...

build:
	go build ./cmd/usdt-rate

test:
	go test -v ./...

docker-build:
	docker build -t usdt-rate .

make run:
	go run ./cmd/usdt-rate/main.go

lint:
	golangci-lint run --fix