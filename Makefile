# Copyright 2022 envd Authors.
#
# The old school Makefile, following are required targets. The Makefile is written
# to allow building multiple binaries. You are free to add more targets or change
# existing implementations, as long as the semantics are preserved.
#
#   make              - default to 'build-local' target
#   make lint         - code analysis
#   make test         - run unit test (or plus integration test)
#   make build-local  - build local binary targets
#   make build-linux  - build linux binary targets
#   make container    - build containers
#   $ docker login registry -u username -p xxxxx
#   make push         - push containers
#   make clean        - clean up targets
#
# Not included but recommended targets:
#   make e2e-test
#
# The makefile is also responsible to populate project version information.
#

#
# Tweak the variables based on your project.
#

# This repo's root import path (under GOPATH).
ROOT := github.com/tensorchord/envd

# Target binaries. You can build multiple binaries for a single project.
TARGETS := envd envd-sshd

# Container image prefix and suffix added to targets.
# The final built images are:
#   $[REGISTRY]/$[IMAGE_PREFIX]$[TARGET]$[IMAGE_SUFFIX]:$[VERSION]
# $[REGISTRY] is an item from $[REGISTRIES], $[TARGET] is an item from $[TARGETS].
IMAGE_PREFIX ?= $(strip )
IMAGE_SUFFIX ?= $(strip )

# Container registries.
REGISTRY ?= ghcr.io/tensorchord

# Container registry for base images.
BASE_REGISTRY ?= docker.io

# Disable CGO by default.
CGO_ENABLED ?= 0

#
# These variables should not need tweaking.
#

# It's necessary to set this because some environments don't link sh -> bash.
export SHELL := bash

# It's necessary to set the errexit flags for the bash shell.
export SHELLOPTS := errexit

# Project main package location (can be multiple ones).
CMD_DIR := ./cmd

# Project output directory.
OUTPUT_DIR := ./bin
DEBUG_DIR := ./debug-bin

# Build directory.
BUILD_DIR := ./build

# Current version of the project.
VERSION ?= $(shell git describe --match 'v[0-9]*' --always --tags --abbrev=0)
BUILD_DATE=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_TAG ?= $(shell if [ -z "`git status --porcelain`" ]; then git describe --exact-match --tags HEAD 2>/dev/null; fi)
GIT_TREE_STATE=$(shell if [ -z "`git status --porcelain`" ]; then echo "clean" ; else echo "dirty"; fi)
GITSHA ?= $(shell git rev-parse --short HEAD)
GIT_LATEST_TAG ?= $(shell git describe --tags --abbrev=0)

# Track code version with Docker Label.
DOCKER_LABELS ?= git-describe="$(shell date -u +v%Y%m%d)-$(shell git describe --tags --always --dirty)"

# Golang standard bin directory.
GOPATH ?= $(shell go env GOPATH)
GOROOT ?= $(shell go env GOROOT)
BIN_DIR := $(GOPATH)/bin
GOLANGCI_LINT := $(BIN_DIR)/golangci-lint

# Default golang flags used in build and test
# -mod=vendor: force go to use the vendor files instead of using the `$GOPATH/pkg/mod`
# -p: the number of programs that can be run in parallel
# -count: run each test and benchmark 1 times. Set this flag to disable test cache
export GOFLAGS ?= -count=1

#
# Define all targets. At least the following commands are required:
#

# All targets.
.PHONY: help lint test build dev container push addlicense debug debug-local build-local generate clean test-local addlicense-install mockgen-install pypi-build base-image envd-lint envd-fmt

.DEFAULT_GOAL:=build-local

build-release:
	@for target in $(TARGETS); do                                                      \
	  CGO_ENABLED=$(CGO_ENABLED) go build -trimpath -v -o $(OUTPUT_DIR)/$${target}     \
	    -ldflags "-s -w -X $(ROOT)/pkg/version.version=$(VERSION) \
		-X $(ROOT)/pkg/version.buildDate=$(BUILD_DATE) \
		-X $(ROOT)/pkg/version.gitCommit=$(GIT_COMMIT) \
		-X $(ROOT)/pkg/version.gitTreeState=$(GIT_TREE_STATE)                     \
		-X $(ROOT)/pkg/version.gitTag=$(GIT_TAG)" \
	    $(CMD_DIR)/$${target};                                                         \
	done
	@$(MAKE) generate-git-tag-info

generate-git-tag-info:
	[[ ! -z "$(GIT_TAG)" ]] && echo "$(GIT_TAG)" > .GIT_TAG_INFO || true

help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

debug: debug-local  ## Build the debug version of envd

# more info about `GOGC` env: https://github.com/golangci/golangci-lint#memory-usage-of-golangci-lint
lint: $(GOLANGCI_LINT)  ## Lint GO code
	@$(GOLANGCI_LINT) run

$(GOLANGCI_LINT):
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin

mockgen-install:
	go install github.com/golang/mock/mockgen@v1.6.0

addlicense-install:
	go install github.com/google/addlicense@latest

build-local:
	@for target in $(TARGETS); do                                                      \
	  CGO_ENABLED=$(CGO_ENABLED) go build -trimpath -v -o $(OUTPUT_DIR)/$${target}     \
	    -ldflags "-s -w -X $(ROOT)/pkg/version.version=$(VERSION) \
		-X $(ROOT)/pkg/version.buildDate=$(BUILD_DATE) \
		-X $(ROOT)/pkg/version.gitCommit=$(GIT_COMMIT) \
		-X $(ROOT)/pkg/version.gitTreeState=$(GIT_TREE_STATE)                     \
		-X $(ROOT)/pkg/version.gitTag=$(GIT_LATEST_TAG) \
		-X $(ROOT)/pkg/version.developmentFlag=true" \
	    $(CMD_DIR)/$${target};                                                         \
	done

pypi-build: clean
	@python3 setup.py sdist bdist_wheel

dev: clean build-local  ## install envd command for local debug
	@python3 setup.py bdist_wheel
	@pip3 install --force-reinstall dist/*.whl

generate: mockgen-install  ## Generate mocks
	@mockgen -source pkg/buildkitd/buildkitd.go -destination pkg/buildkitd/mock/mock.go -package mock
	@mockgen -source pkg/lang/frontend/starlark/interpreter.go -destination pkg/lang/frontend/starlark/mock/mock.go -package mock
	@mockgen -source pkg/progress/compileui/display.go -destination pkg/progress/compileui/mock/mock.go -package mock

# It is used by vscode to attach into the process.
debug-local:
	@for target in $(TARGETS); do                                                      \
	  CGO_ENABLED=$(CGO_ENABLED) go build                                              \
	  	-v -o $(DEBUG_DIR)/$${target}                                                  \
	  	-gcflags='all=-N -l'                                                           \
	    $(CMD_DIR)/$${target};                                                         \
	done

addlicense: addlicense-install  ## Add license to GO code files
	addlicense -c "The envd Authors" $$(find . -type f -name '*.go')

test-local:
	@go test -v -race -coverprofile=coverage.out ./...

test: generate  ## Run the tests
	@go test -race -coverpkg=./pkg/... -coverprofile=coverage.out $(shell go list ./... | grep -v e2e)
	@go tool cover -func coverage.out | tail -n 1 | awk '{ print "Total coverage: " $$3 }'

e2e-test:
	@go test -ldflags "-s -w -X $(ROOT)/pkg/version.version=$(VERSION) \
		-X $(ROOT)/pkg/version.buildDate=$(BUILD_DATE) \
		-X $(ROOT)/pkg/version.gitCommit=$(GIT_COMMIT) \
		-X $(ROOT)/pkg/version.gitTreeState=$(GIT_TREE_STATE)                     \
		-X $(ROOT)/pkg/version.gitTag="$(shell git describe --tags --abbrev=0)" \
		-X $(ROOT)/pkg/version.developmentFlag=true" \
		-race -v -timeout 20m -coverpkg=./pkg/... -coverprofile=e2e-coverage.out ./e2e


e2e-cli-test:
	@go test -ldflags "-s -w -X $(ROOT)/pkg/version.version=$(VERSION) \
		-X $(ROOT)/pkg/version.buildDate=$(BUILD_DATE) \
		-X $(ROOT)/pkg/version.gitCommit=$(GIT_COMMIT) \
		-X $(ROOT)/pkg/version.gitTreeState=$(GIT_TREE_STATE)                     \
		-X $(ROOT)/pkg/version.gitTag="$(shell git describe --tags --abbrev=0)" \
		-X $(ROOT)/pkg/version.developmentFlag=true" \
		-race -v -timeout 20m -coverpkg=./pkg/... -coverprofile=e2e-cli-coverage.out ./e2e/v0/cli

e2e-lang-test:
	@go test -ldflags "-s -w -X $(ROOT)/pkg/version.version=$(VERSION) \
		-X $(ROOT)/pkg/version.buildDate=$(BUILD_DATE) \
		-X $(ROOT)/pkg/version.gitCommit=$(GIT_COMMIT) \
		-X $(ROOT)/pkg/version.gitTreeState=$(GIT_TREE_STATE)                     \
		-X $(ROOT)/pkg/version.gitTag="$(shell git describe --tags --abbrev=0)" \
		-X $(ROOT)/pkg/version.developmentFlag=true" \
		-race -v -timeout 20m -coverpkg=./pkg/... -coverprofile=e2e-lang-coverage.out ./e2e/v0/language

e2e-doc-test:
	@go test -ldflags "-s -w -X $(ROOT)/pkg/version.version=$(VERSION) \
		-X $(ROOT)/pkg/version.buildDate=$(BUILD_DATE) \
		-X $(ROOT)/pkg/version.gitCommit=$(GIT_COMMIT) \
		-X $(ROOT)/pkg/version.gitTreeState=$(GIT_TREE_STATE)                     \
		-X $(ROOT)/pkg/version.gitTag="$(shell git describe --tags --abbrev=0)" \
		-X $(ROOT)/pkg/version.developmentFlag=true" \
		-race -v -timeout 60m -coverpkg=./pkg/... -coverprofile=e2e-doc-coverage.out ./e2e/v0/docs

e2e-cli-test-v1:
	@go test -ldflags "-s -w -X $(ROOT)/pkg/version.version=$(VERSION) \
		-X $(ROOT)/pkg/version.buildDate=$(BUILD_DATE) \
		-X $(ROOT)/pkg/version.gitCommit=$(GIT_COMMIT) \
		-X $(ROOT)/pkg/version.gitTreeState=$(GIT_TREE_STATE)                     \
		-X $(ROOT)/pkg/version.gitTag="$(shell git describe --tags --abbrev=0)" \
		-X $(ROOT)/pkg/version.developmentFlag=true" \
		-race -v -timeout 20m -coverpkg=./pkg/... -coverprofile=e2e-cli-v1-coverage.out ./e2e/v1/cli

e2e-lang-test-v1:
	@go test -ldflags "-s -w -X $(ROOT)/pkg/version.version=$(VERSION) \
		-X $(ROOT)/pkg/version.buildDate=$(BUILD_DATE) \
		-X $(ROOT)/pkg/version.gitCommit=$(GIT_COMMIT) \
		-X $(ROOT)/pkg/version.gitTreeState=$(GIT_TREE_STATE)                     \
		-X $(ROOT)/pkg/version.gitTag="$(shell git describe --tags --abbrev=0)" \
		-X $(ROOT)/pkg/version.developmentFlag=true" \
		-race -v -timeout 20m -coverpkg=./pkg/... -coverprofile=e2e-lang-v1-coverage.out ./e2e/v1/language

e2e-doc-test-v1:
	@go test -ldflags "-s -w -X $(ROOT)/pkg/version.version=$(VERSION) \
		-X $(ROOT)/pkg/version.buildDate=$(BUILD_DATE) \
		-X $(ROOT)/pkg/version.gitCommit=$(GIT_COMMIT) \
		-X $(ROOT)/pkg/version.gitTreeState=$(GIT_TREE_STATE)                     \
		-X $(ROOT)/pkg/version.gitTag="$(shell git describe --tags --abbrev=0)" \
		-X $(ROOT)/pkg/version.developmentFlag=true" \
		-race -v -timeout 60m -coverpkg=./pkg/... -coverprofile=e2e-doc-v1-coverage.out ./e2e/v1/docs

clean:  ## Clean the outputs and artifacts
	@-rm -vrf ${OUTPUT_DIR}
	@-rm -vrf ${DEBUG_DIR}
	@-rm -vrf build dist .eggs *.egg-info

fmt: ## Run go fmt against code.
	go fmt ./...

vet: ## Run go vet against code.
	go vet ./...

black-install:
	@pip install -q black[jupyter]

envd-lint: black-install
	black --check --include '(\.envd|\.py|\.ipynb)$$' .

envd-fmt: black-install
	black --include '(\.envd|\.py|\.ipynb)$$' .
