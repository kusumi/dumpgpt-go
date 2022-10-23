all:
	go build -o dumpgpt-go
fmt:
	go fmt
test:
	go test -v
clean:
	go clean

lint:
	golangci-lint run
