all:
	$(MAKE) -C docs

VERSION := $(shell git describe --tags --always --dirty)
LDFLAGS := -ldflags "-X 'kask/cmd/kask/commands/version.Version=$(VERSION)'"

.PHONY: build test install

build:
	@echo "Version $(VERSION)..."
	mkdir -p "build/$(VERSION)"
	GOOS=darwin  GOARCH=amd64 go build -trimpath $(LDFLAGS) -o build/$(VERSION)/kask-$(VERSION)-darwin-amd64  ./cmd/kask
	GOOS=darwin  GOARCH=arm64 go build -trimpath $(LDFLAGS) -o build/$(VERSION)/kask-$(VERSION)-darwin-arm64  ./cmd/kask
	GOOS=linux   GOARCH=amd64 go build -trimpath $(LDFLAGS) -o build/$(VERSION)/kask-$(VERSION)-linux-amd64   ./cmd/kask
	GOOS=linux   GOARCH=386   go build -trimpath $(LDFLAGS) -o build/$(VERSION)/kask-$(VERSION)-linux-386     ./cmd/kask
	GOOS=linux   GOARCH=arm   go build -trimpath $(LDFLAGS) -o build/$(VERSION)/kask-$(VERSION)-linux-arm     ./cmd/kask
	GOOS=linux   GOARCH=arm64 go build -trimpath $(LDFLAGS) -o build/$(VERSION)/kask-$(VERSION)-linux-arm64   ./cmd/kask
	GOOS=freebsd GOARCH=amd64 go build -trimpath $(LDFLAGS) -o build/$(VERSION)/kask-$(VERSION)-freebsd-amd64 ./cmd/kask
	GOOS=freebsd GOARCH=386   go build -trimpath $(LDFLAGS) -o build/$(VERSION)/kask-$(VERSION)-freebsd-386   ./cmd/kask
	GOOS=freebsd GOARCH=arm   go build -trimpath $(LDFLAGS) -o build/$(VERSION)/kask-$(VERSION)-freebsd-arm   ./cmd/kask

.PHONY: install

install:
	go build $(LDFLAGS) -o ~/bin/kask  ./cmd/kask
