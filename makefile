APP := totion

ifeq ($(OS),Windows_NT)
BIN := $(APP).exe
else
BIN := $(APP)
endif

build:
	@go build -o $(BIN) .

run: build
	@./$(BIN)