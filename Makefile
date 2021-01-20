#$$ Basic Makefile for Golang project
SERVICE			?= $(shell basename `go list`)
PACKAGE			?= $(shell go list)
PACKAGES		?= $(shell go list ./...)
FILES			?= $(shell find . -type f -name '*.go' -not -path "./vendor/*")
TAG         	?= $(shell git describe --tags --abbrev=0)
NOW				?= $(shell date +'%d-%m-%Y_%T')
COMMIT      	?= $(shell git rev-parse HEAD)
CURRENT_DIR     ?= $(shell pwd)

#$$ Binaries
GOCMD			:=	go
GOBUILD			:=	$(GOCMD) build
GOMOD			:=	$(GOCMD) mod
GOTEST			:=	$(GOCMD) test
GOFMT			:=	$(GOCMD) fmt
GOCLEAN			:=	$(GOCMD) clean

#$$
PROTOC			?= protoc

.PHONY: help clean fmt lint vet test coverage build all

default: help

help: ## show this help
	@echo 'usage: make [target] ...'
	@echo ''
	@echo 'targets:'
	@egrep '^(.+)\:\ .*##\ (.+)' ${MAKEFILE_LIST} | sed 's/:.*##/#/' | column -t -c 2 -s '#'
	@echo ''

all: ## clean, format, build and unit test
	@make clean
	@make fmt
	@make test
	@make build

install: ## build and install go application executable
	@go install -v ./...

env: ## Print useful environment variables to stdout
	@echo $(CURDIR)
	@echo $(SERVICE)
	@echo $(PACKAGE)

clean: ## clean build and coverage results
	@$(GOCLEAN) -i ./...
	@rm -rf $(BINARY) coverage

tools: ## fetch and install all required tools
	@go get -u golang.org/x/tools/cmd/goimports
	@go get -u github.com/golang/protobuf/protoc-gen-go
	@go get -u github.com/mwitkow/go-proto-validators/protoc-gen-govalidators
	@go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	@go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
	@go get -u github.com/matryer/moq
	@go get -u github.com/axw/gocov/gocov
	@go get -u github.com/AlekSi/gocov-xml

fmt: ## format the go source files
	@$(GOFMT) ./...
	@goimports -w $(FILES)

lint: ## run go lint on the source files
	@golint $(PACKAGES)

vet: ## run go vet on the source files
	@go vet ./...

doc: ## generate godocs and start a local documentation webserver on port 8085
	@godoc -http=:8085 -index

update-dependencies: ## update golang dependencies
	@$(GOMOD) download

test-all: test test-bench coverage

test: ## Run short tests
	@$(GOTEST) -v ./... -short

test-it: ## Run all tests
	@$(GOTEST) -v ./...

test-bench: ## run benchmark tests
	@$(GOTEST) -bench ./...

# Generate test coverage
coverage: ## Run test coverage and generate HTML/XML report
	@rm -fr coverage
	@mkdir coverage
	@go list -f '{{if gt (len .TestGoFiles) 0}}"go test -covermode count -coverprofile {{.Name}}.coverprofile -coverpkg ./... {{.ImportPath}}"{{end}}' ./... | xargs -I {} bash -c {}
	@echo "mode: count" > coverage/cover.out
	@grep -h -v "^mode:" *.coverprofile >> "coverage/cover.out"
	@rm -f *.coverprofile
	@go tool cover -html=coverage/cover.out -o=coverage/coverage.html
	@gocov convert coverage/cover.out | gocov-xml > coverage/coverage.xml

version: build
	@./$(BINARY) version --json

build: ## generate binary
	@CGO_ENABLED=0 GOOS=linux $(GOBUILD) -a -installsuffix cgo -ldflags="${BUILD_NOW} ${BUILD_VERSION} ${BUILD_COMMIT}" -o $(BINARY)
	@strip $(BINARY)
