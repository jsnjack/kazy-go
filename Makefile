BINARY:=kazy
PWD:=$(shell pwd)
BUILD_TYPES:=rpm deb
VERSION=0.0.0
MONOVA:=$(shell which monova dot 2> /dev/null)

version:
ifdef MONOVA
override VERSION="$(shell monova)"
else
	$(info "Install monova (https://github.com/jsnjack/monova) to calculate version")
endif

test:
	go test

coverage:
	go test -coverprofile .coverage && go tool cover -html=.coverage && go tool cover -html=.coverage

build: version
	go build -ldflags="-X main.version=${VERSION}" -o ${BINARY}

dist: build
	@for type in ${BUILD_TYPES} ; do \
		cd ${PWD}/dist && fpm --input-type dir --output-type $$type \
		--name kazy-go --version ${VERSION} --license MIT --no-depends --provides kazy \
		--vendor jsnjack@gmail.com \
		--maintainer jsnjack@gmail.com --description "Highlights output from STDIN" \
		--url https://github.com/jsnjack/kazy-go --force --chdir ${PWD} ./kazy=/usr/bin/kazy; \
	done

clean:
	go clean

.PHONY: dist
