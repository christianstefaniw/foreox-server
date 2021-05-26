SRC_DIR = src
BIN_DIR = bin
BIN_NAME = devcord
DEVCORD = cmd/devcord/main.go

build:
	go build -o $(BIN_DIR)/$(BIN_NAME) $(SRC_DIR)/$(DEVCORD)

run:
	go run $(SRC_DIR)/$(DEVCORD)