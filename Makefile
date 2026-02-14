BINARY = bilingual_pdf
VERSION = $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS = -ldflags "-X bilingual_pdf/cmd.Version=$(VERSION)"

.PHONY: build test test-integration test-cover lint clean install \
	dist dist-darwin-arm64 dist-darwin-amd64 dist-linux-amd64 dist-windows-amd64

build:
	go build $(LDFLAGS) -o $(BINARY) .

install:
	go install $(LDFLAGS) .

test:
	go test ./...

test-integration:
	go test -tags=integration ./...

test-cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run

DIST_DIR = dist

dist: dist-darwin-arm64 dist-darwin-amd64 dist-linux-amd64 dist-windows-amd64

dist-darwin-arm64:
	mkdir -p $(DIST_DIR)
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY) .
	tar -czf $(DIST_DIR)/$(BINARY)-$(VERSION)-darwin-arm64.tar.gz -C $(DIST_DIR) $(BINARY)
	rm $(DIST_DIR)/$(BINARY)

dist-darwin-amd64:
	mkdir -p $(DIST_DIR)
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY) .
	tar -czf $(DIST_DIR)/$(BINARY)-$(VERSION)-darwin-amd64.tar.gz -C $(DIST_DIR) $(BINARY)
	rm $(DIST_DIR)/$(BINARY)

dist-linux-amd64:
	mkdir -p $(DIST_DIR)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY) .
	tar -czf $(DIST_DIR)/$(BINARY)-$(VERSION)-linux-amd64.tar.gz -C $(DIST_DIR) $(BINARY)
	rm $(DIST_DIR)/$(BINARY)

dist-windows-amd64:
	mkdir -p $(DIST_DIR)
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY).exe .
	zip -j $(DIST_DIR)/$(BINARY)-$(VERSION)-windows-amd64.zip $(DIST_DIR)/$(BINARY).exe
	rm $(DIST_DIR)/$(BINARY).exe

clean:
	rm -f $(BINARY)
	rm -f coverage.out coverage.html
	rm -f testdata/*.pdf testdata/*.html
	rm -rf $(DIST_DIR)
