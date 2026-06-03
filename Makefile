BINARY := bin/omapctl
GO := go

.PHONY: build lint test clean

build:
	CGO_ENABLED=1 $(GO) build -o $(BINARY) .

test:
	$(GO) test ./internal/output/...

lint:
	golangci-lint run ./...

clean:
	rm -f $(BINARY)
