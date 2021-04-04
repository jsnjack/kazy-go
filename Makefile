PWD:=$(shell pwd)
VERSION=0.0.0
MONOVA:=$(shell which monova dot 2> /dev/null)

version:
ifdef MONOVA
override VERSION=$(shell monova)
else
	$(info "Install monova (https://github.com/jsnjack/monova) to calculate version")
endif

bin/kazy: bin/kazy_linux_amd64
	cp bin/kazy_linux_amd64 bin/kazy

bin/kazy_linux_amd64: version main.go cmd/*.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-X github.com/jsnjack/kazy-go/cmd.Version=${VERSION}" -o bin/kazy_linux_amd64

bin/kazy_darwin_amd64: version main.go cmd/*.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-X github.com/jsnjack/kazy-go/cmd.Version=${VERSION}" -o bin/kazy_darwin_amd64

build: bin/kazy bin/kazy_linux_amd64 bin/kazy_darwin_amd64

release: build
	tar --transform='s,_.*,,' --transform='s,bin/,,' -cz -f bin/kazy_linux_amd64.tar.gz bin/kazy_linux_amd64
	tar --transform='s,_.*,,' --transform='s,bin/,,' -cz -f bin/kazy_darwin_amd64.tar.gz bin/kazy_darwin_amd64
	grm release jsnjack/kazy-go -f bin/kazy_linux_amd64.tar.gz -f bin/kazy_darwin_amd64.tar.gz -t "v`monova`"

test:
	go test github.com/jsnjack/kazy-go/cmd

.PHONY: version release build test
