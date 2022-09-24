.PHONY: *.go go.mod go.sum

build:
	CGO_ENABLED=0 go build -o bin/cvault ./cmd/cvault
