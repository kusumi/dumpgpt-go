all:
	go build -o dumpgpt
fmt:
	go fmt
test:
	go test -v
clean:
	go clean
	rm ./dumpgpt
