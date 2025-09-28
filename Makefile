.DEFAULT_GOAL := all

license:
	curl -sL https://liam.sh/-/gh/g/license-header.sh | bash -s

up:
	go get -u ./... && go mod tidy
	go get -u -t ./... && go mod tidy

generate: license
	go generate -x ./...
	cd examples/simple && go run . generate-markdown > README.md

test:
	gofmt -e -s -w .
	go vet .
	go test -v ./...

all: generate test
	@echo
