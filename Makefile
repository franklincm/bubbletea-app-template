
.DEFAULT_GOAL := build

fmt:
	gofmt -w -s .
.PHONY:fmt

lint: fmt
	golint .
.PHONY:lint

vet: fmt
	go vet ./...
.PHONY:vet

build: vet
	go build -o dist/template *.go
.PHONY:build

test:
	go test
.PHONY:test

clean:
	$(RM) -rf dist/
.PHONY:clean
