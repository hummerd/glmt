
default: build


build: VERSION = $(shell git describe --tags)
build: VERSION_DATE = $(shell git show -s --format=%ct HEAD)
build: PACKAGE = $(shell go list -m)
build: LDFLAGS = --ldflags "-X '$(PACKAGE)/internal/glmt.Version=$(VERSION)' -X '$(PACKAGE)/internal/glmt.VersionDate=$(VERSION_DATE)'"
build:
	go build $(LDFLAGS) ./cmd/glmt

test:
	go test -race ./...

lint:
	golangci-lint run

clean:
	rm glmt
