GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOLIST=$(GOCMD) list
GODEP=godep save
BINARY_NAME=falcon-judge

GIT_COMMIT=$(shell git rev-parse HEAD)
BUILD_TIME=$(shell date "+%Y-%m-%dT%H:%M:%S%z")
LDFLAGS=-ldflags="-X main.GitCommit=$(GIT_COMMIT) -X main.BuildTime=$(BUILD_TIME)"

test:
		$(GOTEST) -v ./...
clean:
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
run:
		$(GOBUILD) -o $(BINARY_NAME) -v ./...
		./$(BINARY_NAME)
deps:
		$(GOGET) github.com/tools/godep
		$(GODEP) -v ./...


# Cross compilation
build:
		$(GOBUILD) -o $(BINARY_NAME) "$(LDFLAGS)" -v
build-linux:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) -v
build-windows:
		CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME) -v

pack: build-linux
	@echo "==> Packing ..."
	@tar czvf $(BINARY_NAME)-linux-amd64.tar.gz $(BINARY_NAME)
	@echo ""
	@rm $(BINARY_NAME)
