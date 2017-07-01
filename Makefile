BINARY:=kazy
PWD:=$(shell pwd)
BUILD_TYPES:=rpm deb
VERSION:=$(shell ./${BINARY} --version)


test:
	go test

coverage:
	go test -coverprofile .coverage && go tool cover -html=.coverage && go tool cover -html=.coverage

build:
	go build -o ${BINARY}

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
