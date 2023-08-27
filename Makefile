
.DEFAULT_GOAL := build

precommit:
	pre-commit install
.PHONY:precommit

fmt: precommit
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

vhs:
	vhs demo.tape
.PHONY: vhs

clean:
	$(RM) -rf dist/
.PHONY:clean
