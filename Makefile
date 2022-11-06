BIN_DIR=./bin

help: Makefile
	@echo
	@echo " Choose a command run:"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'

## dev: run server
dev:
	go build -o $(BIN_DIR)/server -v ./cmd/server
	$(BIN_DIR)/server

## clean: clean bin directory
clean:
	rm $(BIN_DIR)/*

## lint: run golangci-lint
lint:
	golangci-lint run ./...

## mod_tidy: go mod tidy
mod_tidy:
	go mod tidy

## build: build executable
build:
	go build -o $(BIN_DIR)/server -v ./cmd/server