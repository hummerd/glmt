
default: build

build:
	go build -o glmt ./cmd

test:
	go test -race ./...

lint:
	golangci-lint run

clean:
	rm glmt
