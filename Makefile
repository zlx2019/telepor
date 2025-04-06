APP_NAME := telepor
OUT_DIR := bin
FLAGS := CGO_ENABLED=0
BUILD_COMMANDS := go build -ldflags="-s -w"

build:
	@$(BUILD_COMMANDS) -o $(OUT_DIR)/$(APP_NAME)

build-linux:
	@$(FLAGS) GOOS=linux GOARCH=amd64 $(BUILD_COMMANDS) -o $(OUT_DIR)/linux/$(APP_NAME)

build-windows:
	@$(FLAGS) GOOS=windows GOARCH=amd64 $(BUILD_COMMANDS) -o $(OUT_DIR)/windows/$(APP_NAME)

build-all: build build-linux build-windows

clean:
	@rm -rf $(OUT_DIR)