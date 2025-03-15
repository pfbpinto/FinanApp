run:
	go run cmd/app/main.go

build:
	go build -o finanapp cmd/app/main.go

test:
	go test ./...

lint:
	golangci-lint run
