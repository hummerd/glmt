
default: build

build:
	go build ./cmd/glmt

test:
	go test -race ./...

lint:
	golangci-lint run

clean:
	rm glmt
