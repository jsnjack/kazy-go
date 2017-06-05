BINARY=kazy
GITHUB_ACCESS_TOKEN?=

test:
	go test

build:
	go build -o ${BINARY}

clean:
	go clean

install: clean build
	cp ./${BINARY} /usr/bin/

uninstall:
	rm -f /usr/bin/${BINARY}

compress: clean build
	tar -czf ./dist/${BINARY}-${shell ./${BINARY} --version}.tar.gz kazy

release: test compress
	git tag v${shell ./${BINARY} --version}
	git push --tags
	github-release release --user jsnjack --repo kazy-go --tag v${shell ./${BINARY} --version}
	github-release upload --user jsnjack --repo kazy-go --tag v${shell ./${BINARY} --version} -f ./dist/${BINARY}-${shell ./${BINARY} --version}.tar.gz -n ${BINARY}-${shell ./${BINARY} --version}.tar.gz

.PHONY: build clean install uninstall release
