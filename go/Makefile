# Kanije Kalesi — Makefile
# Tüm hedefler: make help

BINARY     := kanije
MODULE     := github.com/kanije-kalesi/kanije
VERSION    := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DATE := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
LDFLAGS    := -s -w -X main.Version=$(VERSION) -X main.BuildDate=$(BUILD_DATE)

DIST := dist

.DEFAULT_GOAL := help

## help: Bu yardım mesajını göster
.PHONY: help
help:
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## //' | column -t -s ':'

## build: Mevcut platform için derle
.PHONY: build
build:
	go build -ldflags="$(LDFLAGS)" -o $(BINARY) ./cmd/kanije/

## build-windows: Windows AMD64 için cross-compile et (CGo yok)
.PHONY: build-windows
build-windows:
	mkdir -p $(DIST)
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 \
	go build -ldflags="$(LDFLAGS) -H=windowsgui" \
	-o $(DIST)/$(BINARY)-windows-amd64.exe ./cmd/kanije/

## build-linux: Linux AMD64 için cross-compile et
.PHONY: build-linux
build-linux:
	mkdir -p $(DIST)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
	go build -ldflags="$(LDFLAGS)" \
	-o $(DIST)/$(BINARY)-linux-amd64 ./cmd/kanije/

## build-arm64: Raspberry Pi 4/5 (Linux ARM64) için cross-compile et
.PHONY: build-arm64
build-arm64:
	mkdir -p $(DIST)
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 \
	go build -ldflags="$(LDFLAGS)" \
	-o $(DIST)/$(BINARY)-linux-arm64 ./cmd/kanije/

## build-arm: Raspberry Pi 3 / eski Pi (Linux ARM 32-bit) için cross-compile et
.PHONY: build-arm
build-arm:
	mkdir -p $(DIST)
	GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 \
	go build -ldflags="$(LDFLAGS)" \
	-o $(DIST)/$(BINARY)-linux-arm ./cmd/kanije/

## build-all: Tüm platformlar için derle
.PHONY: build-all
build-all: build-windows build-linux build-arm64 build-arm
	@echo ""
	@echo "Derlenen binary'ler:"
	@ls -lh $(DIST)/

## compress: UPX ile binary boyutunu küçült (UPX kurulu olmalı)
.PHONY: compress
compress: build-all
	@if command -v upx >/dev/null 2>&1; then \
		upx --best $(DIST)/$(BINARY)-*; \
	else \
		echo "UPX bulunamadı, atlıyor."; \
	fi

## run: Yerel binary'yi başlat (debug modu)
.PHONY: run
run:
	KANIJE_LOG_LEVEL=debug go run ./cmd/kanije/ start

## test: Testleri çalıştır
.PHONY: test
test:
	go test -v -race ./...

## lint: Kod kalitesini kontrol et
.PHONY: lint
lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint bulunamadı: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

## tidy: go.sum'u güncelle ve kullanılmayan bağımlılıkları temizle
.PHONY: tidy
tidy:
	go mod tidy

## install-windows: Windows'a kur (Task Scheduler)
.PHONY: install-windows
install-windows: build-windows
	powershell -ExecutionPolicy Bypass -File deploy\windows\install.ps1

## install-linux: Linux'a kur (systemd)
.PHONY: install-linux
install-linux: build-linux
	sudo bash deploy/linux/install.sh

## install-pi: Raspberry Pi'ya kur (ARM64)
.PHONY: install-pi
install-pi: build-arm64
	sudo bash deploy/linux/install.sh

## clean: Derleme çıktılarını temizle
.PHONY: clean
clean:
	rm -rf $(DIST) $(BINARY) $(BINARY).exe

## version: Versiyon bilgisini göster
.PHONY: version
version:
	@echo "Versiyon   : $(VERSION)"
	@echo "Build Tarihi: $(BUILD_DATE)"
	@echo "Go Versiyonu: $(shell go version)"
