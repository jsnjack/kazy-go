BINARY=kazy
GITHUB_ACCESS_TOKEN?=

test:
	go test

coverage:
	go test -coverprofile .coverage && go tool cover -html=.coverage && go tool cover -html=.coverage

build:
	go build -o ${BINARY}

clean:
	go clean

install: clean build
	cp ./${BINARY} /usr/bin/

web-install:
	cd /tmp && curl -L --output /tmp/kazy.tar.gz ${shell curl -s "https://api.github.com/repos/jsnjack/kazy-go/releases/latest" | jq -r '.assets[] | select(.name=="kazy.tar.gz") | .browser_download_url'}
	cd /tmp/ && tar -xzf kazy.tar.gz
	sudo cp /tmp/${BINARY} /usr/bin/

uninstall:
	rm -f /usr/bin/${BINARY}

compress: clean build
	tar -czf ./dist/${BINARY}-${shell ./${BINARY} --version}.tar.gz kazy

release: test compress
	git tag v${shell ./${BINARY} --version}
	git push --tags
	github-release release --user jsnjack --repo kazy-go --tag v${shell ./${BINARY} --version}
	github-release upload --user jsnjack --repo kazy-go --tag v${shell ./${BINARY} --version} -f ./dist/${BINARY}-${shell ./${BINARY} --version}.tar.gz -n ${BINARY}.tar.gz

.PHONY: build clean install uninstall release
