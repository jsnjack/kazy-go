BINARY  := kazy
REPO    := jsnjack/kazy-go
PKG     := ./...
VERSION := 0.0.0
MONOVA  := $(shell which monova 2> /dev/null)
LDFLAGS  = -ldflags="-X github.com/jsnjack/kazy-go/cmd.Version=$(VERSION)"

export PATH := $(PATH):$(shell go env GOPATH)/bin

version:
ifdef MONOVA
override VERSION = $(shell monova)
override LDFLAGS = -ldflags="-X github.com/jsnjack/kazy-go/cmd.Version=$(VERSION)"
else
	$(info "Install monova with: grm install jsnjack/monova")
endif

test:
	go test $(PKG)

vet:
	go vet $(PKG)

fmt:
	@command -v goimports >/dev/null 2>&1 || { \
	  echo "goimports is not installed. Install it with:"; \
	  echo "  go install golang.org/x/tools/cmd/goimports@latest"; \
	  exit 1; \
	}
	goimports -w .

lint: vet
	@command -v golangci-lint >/dev/null 2>&1 || { \
	  echo "golangci-lint is not installed. Install it with:"; \
	  echo "  grm install golangci/golangci-lint"; \
	  exit 1; \
	}
	golangci-lint run

check: fmt vet build test lint
	@echo "==> make check: all green"

standards:
	curl -sL https://raw.githubusercontent.com/jsnjack/standards/master/AGENTS.universal.md \
	    -o AGENTS.universal.md
	curl -sL https://raw.githubusercontent.com/jsnjack/standards/master/AGENTS.go.md \
	    -o AGENTS.go.md

bin/$(BINARY): bin/$(BINARY)_linux_amd64
	cp bin/$(BINARY)_linux_amd64 bin/$(BINARY)
	ln -sf bin/$(BINARY) $(BINARY)
bin/$(BINARY)_linux_amd64: version
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY)_linux_amd64
bin/$(BINARY)_linux_arm64: version
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o bin/$(BINARY)_linux_arm64
bin/$(BINARY)_darwin_amd64: version
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY)_darwin_amd64
bin/$(BINARY)_darwin_arm64: version
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/$(BINARY)_darwin_arm64

build: bin/$(BINARY) bin/$(BINARY)_linux_amd64 bin/$(BINARY)_linux_arm64 bin/$(BINARY)_darwin_amd64 bin/$(BINARY)_darwin_arm64

release: build
	tar -czf bin/$(BINARY)_linux_amd64.tar.gz  --transform 's|.*/$(BINARY)_.*|$(BINARY)|' bin/$(BINARY)_linux_amd64
	tar -czf bin/$(BINARY)_linux_arm64.tar.gz  --transform 's|.*/$(BINARY)_.*|$(BINARY)|' bin/$(BINARY)_linux_arm64
	tar -czf bin/$(BINARY)_darwin_amd64.tar.gz --transform 's|.*/$(BINARY)_.*|$(BINARY)|' bin/$(BINARY)_darwin_amd64
	tar -czf bin/$(BINARY)_darwin_arm64.tar.gz --transform 's|.*/$(BINARY)_.*|$(BINARY)|' bin/$(BINARY)_darwin_arm64
	grm release $(REPO) \
		-f bin/$(BINARY)_linux_amd64.tar.gz \
		-f bin/$(BINARY)_linux_arm64.tar.gz \
		-f bin/$(BINARY)_darwin_amd64.tar.gz \
		-f bin/$(BINARY)_darwin_arm64.tar.gz \
		-t "v`monova`"

clean:
	rm -rf bin/ $(BINARY)

.PHONY: version build release test vet fmt lint check standards clean
