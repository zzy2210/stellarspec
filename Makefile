# Go参数
GO := go
GOFLAGS := $(GOFLAGS) -v
LDFLAGS := -s -w # 缩小二进制文件体积

# 项目配置
BINARY_NAME :=  stellar 
CMD_PATH := ./cmd/stellarspec.go

# 构建目标
.PHONY: all build clean run help

all: build

build:
	@echo "Building $(BINARY_NAME)..."
	$(GO) build $(GOFLAGS) -ldflags="$(LDFLAGS)" -o $(BINARY_NAME) $(CMD_PATH)
	@echo "$(BINARY_NAME) built successfully."

clean:
	@echo "Cleaning..."
	$(GO) clean
	rm -f $(BINARY_NAME)
	@echo "Cleaned."

run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BINARY_NAME) --help

help:
	@echo "Available targets:"
	@echo "  all         - Build the application (default)"
	@echo "  build       - Build the application"
	@echo "  clean       - Remove build artifacts"
	@echo "  run         - Build and run the application with --help"
	@echo "  help        - Show this help message"

# 默认目标
.DEFAULT_GOAL := help