GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o bin/$(GOOS)_$(GOARCH)_app main.go

clean:
	rm -rf ./bin