TYGO := $(shell which tygo 2>/dev/null || echo $(shell go env GOPATH)/bin/tygo)

.PHONY: all build run test gen-types clean

# Default target
all: build

build:
	@mkdir -p bin
	cd daemon && go build -o ../bin/rudder ./cmd/rudder

run:
	cd daemon && go run ./cmd/rudder

test:
	cd daemon && go test -v ./...

gen-types:
	@if [ ! -x "$(TYGO)" ]; then \
		echo "Installing tygo..."; \
		go install github.com/gzuidhof/tygo@latest; \
	fi
	cd daemon && $(TYGO) generate

clean:
	cd daemon && go clean
