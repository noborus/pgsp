BINARY_NAME := pgsp
SRCS := $(shell git ls-files '*.go')
LDFLAGS := "-X github.com/noborus/pgsp/cmd.Version=$(shell git describe --tags --abbrev=0 --always) -X github.com/noborus/pgsp/cmd.Revision=$(shell git rev-parse --verify --short HEAD)"

all: build

test: $(SRCS)
	go test ./...

deps:
	go mod tidy

build: deps $(BINARY_NAME)

$(BINARY_NAME): $(SRCS)
	go build -ldflags $(LDFLAGS) ./cmd/pgsp

install:
	go install -ldflags $(LDFLAGS)

clean:
	rm -f $(BINARY_NAME)

.PHONY: all test deps build install clean
