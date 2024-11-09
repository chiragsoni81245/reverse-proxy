build:
	@go build -o bin/reverse-proxy

run: build
	@./bin/reverse-proxy

test:
	@go test ./... -v --race

