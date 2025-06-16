# Go参数
GO := go
GOFLAGS := $(GOFLAGS) -v
LDFLAGS := -s -w # 缩小二进制文件体积

# 项目配置
BINARY_NAME :=  stellar 
CMD_PATH := ./cmd/stellarspec.go
BUILD_DIR := build


# 构建目标
.PHONY: all build clean run help

all: build

build:
	@echo "Creating build directory..."
	@mkdir -p $(BUILD_DIR)
	@echo "Building $(BINARY_NAME)..."
	$(GO) build $(GOFLAGS) -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_PATH)
	@echo "$(BINARY_NAME) built successfully."

clean:
	@echo "Cleaning..."
	$(GO) clean
	rm -rf $(BUILD_DIR)
	@echo "Cleaned."

run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME) --help


install: build
    @echo "Installing $(BINARY_NAME) to /usr/local/bin..."
    sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
    @echo "$(BINARY_NAME) installed successfully."


help:
	@echo "Available targets:"
	@echo "  all         - Build the application (default)"
	@echo "  build       - Build the application"
	@echo "  clean       - Remove build artifacts"
	@echo "  run         - Build and run the application with --help"
	@echo "  help        - Show this help message"

# 默认目标
.DEFAULT_GOAL := help