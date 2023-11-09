##################################################
# Variables                                      #
##################################################

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

VERSION		   ?= latest

ARCH       ?=amd64
CGO        ?=0
TARGET_OS  ?=linux

GIT_VERSION = $(shell git describe --always --abbrev=7)
GIT_COMMIT  = $(shell git rev-list -1 HEAD)
DATE        = $(shell date -u +"%Y.%m.%d.%H.%M.%S")
GOPATH      = $(shell go env GOPATH)
GOROOT      = $(shell go env GOROOT)
PROTOCPATH  = "${GOPATH}/bin"

GIT_VERSION := $(shell git rev-parse HEAD)
CURRENT_TIME := $(shell date "+%F-%T")
ifndef BUILD_VERSION
	BUILD_VERSION := $(GIT_VERSION)-$(CURRENT_TIME)
endif

.PHONY: binaries swagger
.DEFAULT_GOAL := binaries

##################################################
# Build                                          #
##################################################
GO_BUILD_VARS= GO111MODULE=on CGO_ENABLED=$(CGO) GOOS=$(TARGET_OS) GOARCH=$(ARCH) GOROOT=$(GOROOT) GOPATH=$(GOPATH) GOPRIVATE=$(GOPRIVATE)

fmt: ## Run go fmt against code.
	go fmt ./...

vet: ## Run go vet against code.
	go vet ./...

deps:
	go install github.com/golang/mock/mockgen@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest

go-gen: deps
	go generate ./...

binaries: fmt vet
	$(GO_BUILD_VARS) go build -ldflags "-X k9s-autoscaler/pkg/version.Build=$(BUILD_VERSION)" -o bin/k9s-autoscaler ./cmd/k9s-autoscaler/
