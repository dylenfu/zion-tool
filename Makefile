SHELL=/bin/bash

# Go parameters
GOCMD=GO111MODULE=on go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
ENV=$(ONROBOT)

compile:
	@$(GOBUILD) -o build/$(ENV)/zion-tool cmd/main.go

compile-linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o build/$(ENV)/zion-tool-linux cmd/main.go

run:
	@echo test case $(t)
	./build/$(ENV)/zion-tool -config=build/$(ENV)/config.json -t=$(t)

clean: