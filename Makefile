# NPM Security Scanner Makefile

.PHONY: build test clean run install deps lint format help

# デフォルトターゲット
.DEFAULT_GOAL := help

# バイナリ名
BINARY_NAME := npm-security-scanner
BUILD_DIR := ./bin

# Go関連の設定
GO_VERSION := 1.21
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

# ビルドフラグ
LDFLAGS := -ldflags "-s -w -X main.appVersion=$(shell git describe --tags --always --dirty 2>/dev/null || echo 'dev')"

## help: このヘルプメッセージを表示
help:
	@echo "Available commands:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'

## deps: 依存関係をインストール
deps:
	@echo "📦 Installing dependencies..."
	go mod download
	go mod tidy

## build: バイナリをビルド
build: deps
	@echo "🔨 Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "✅ Build completed: $(BUILD_DIR)/$(BINARY_NAME)"

## build-all: 全プラットフォーム向けにビルド
build-all: deps
	@echo "🔨 Building for all platforms..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	@echo "✅ Cross-platform build completed"

## test: テストを実行
test: deps
	@echo "🧪 Running tests..."
	go test -v -race -coverprofile=coverage.out ./...
	@echo "📊 Coverage report:"
	go tool cover -func=coverage.out

## test-coverage: カバレッジレポートをHTMLで表示
test-coverage: test
	@echo "📈 Opening coverage report in browser..."
	go tool cover -html=coverage.out

## bench: ベンチマークテストを実行
bench: deps
	@echo "⚡ Running benchmarks..."
	go test -bench=. -benchmem ./...

## lint: コードの静的解析
lint:
	@echo "🔍 Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint not found. Installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run; \
	fi

## format: コードをフォーマット
format:
	@echo "✨ Formatting code..."
	go fmt ./...
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -w .; \
	else \
		echo "📦 Installing goimports..."; \
		go install golang.org/x/tools/cmd/goimports@latest; \
		goimports -w .; \
	fi

## run: プログラムを実行（引数: DIR）
run: build
	@echo "🚀 Running $(BINARY_NAME)..."
	$(BUILD_DIR)/$(BINARY_NAME) $(DIR)

## install: システムにインストール
install: build
	@echo "📦 Installing to system..."
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "✅ Installed to /usr/local/bin/$(BINARY_NAME)"

## clean: ビルド成果物を削除
clean:
	@echo "🧹 Cleaning up..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out
	go clean

## security: セキュリティチェック
security:
	@echo "🔒 Running security checks..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "📦 Installing gosec..."; \
		go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest; \
		gosec ./...; \
	fi

## mod-update: Go modulesを最新に更新
mod-update:
	@echo "📦 Updating Go modules..."
	go get -u ./...
	go mod tidy

## docker-build: Dockerイメージをビルド
docker-build:
	@echo "🐳 Building Docker image..."
	docker build -t $(BINARY_NAME):latest .

## pre-commit: コミット前のチェック
pre-commit: format lint test security
	@echo "✅ Pre-commit checks completed"
