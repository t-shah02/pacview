.PHONY: run build test

run:
	go run main.go

build:
	go build -o bin/pacview main.go

test:
	go test ./internal/utils_test/...