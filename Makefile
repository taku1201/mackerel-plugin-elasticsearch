BINARY_NAME := mackerel-plugin-elasticsearch
LDFLAGS := "-w -s"

.PHONY: build clean

build:
	GOARCH=arm64 GOOS=darwin go build -o build/${BINARY_NAME}-darwin-arm64 -ldflags=${LDFLAGS} -trimpath ./cmd/${BINARY_NAME}
	GOARCH=arm64 GOOS=linux go build -o build/${BINARY_NAME}-linux-arm64 -ldflags=${LDFLAGS} -trimpath ./cmd/${BINARY_NAME}
	GOARCH=amd64 GOOS=darwin go build -o build/${BINARY_NAME}-darwin-amd64 -ldflags=${LDFLAGS} -trimpath ./cmd/${BINARY_NAME}
	GOARCH=amd64 GOOS=linux go build -o build/${BINARY_NAME}-linux-amd64 -ldflags=${LDFLAGS} -trimpath ./cmd/${BINARY_NAME}

clean:
	rm -f build/${BINARY_NAME}-*

all: clean build ;
