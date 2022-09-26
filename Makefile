.PHONY: *.go go.mod go.sum

native-build:
	CGO_ENABLED=0 go build -o ./bin/cvault ./cmd/cvault

build:
	goreleaser build --snapshot --rm-dist

release:
	goreleaser release --snapshot --rm-dist
