.PHONY: *.go go.mod go.sum

build:
	goreleaser build --snapshot --rm-dist

release:
	goreleaser release --snapshot --rm-dist
