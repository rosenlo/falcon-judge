GIT_COMMIT=$(shell git rev-parse HEAD)
BUILD_TIME=$(shell date "+%Y-%m-%dT%H:%M:%S%z")
VERSION_FLAG=-ldflags="-X main.GitCommit=$(GIT_COMMIT) -X main.BuildTime=$(BUILD_TIME)"

all: b

b:
	GOOS=linux GOARCH=amd64 go build $(VERSION_FLAG)

install:
	GOOS=linux GOARCH=amd64 go install $(VERSION_FLAG) .
