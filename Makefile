GOOS := linux
GOARCH := amd64
BINARY_NAME := format-postman

build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "-s -w" -o $(BINARY_NAME)
