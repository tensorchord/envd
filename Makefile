# Copyright 2022 MIDI Authors.
#
# The old school Makefile, following are required targets. The Makefile is written
# to allow building multiple binaries. You are free to add more targets or change
# existing implementations, as long as the semantics are preserved.
#
#   make              - default to 'build' target
#   make lint         - code analysis
#   make test         - run unit test (or plus integration test)
#   make build        - alias to build-local target
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
ROOT := github.com/tensorchord/midi

# Target binaries. You can build multiple binaries for a single project.
TARGETS := midi

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

# Enable CGO by default.
CGO_ENABLED ?= 1

#
# These variables should not need tweaking.
#

# It's necessary to set this because some environments don't link sh -> bash.
export SHELL := /bin/bash

# It's necessary to set the errexit flags for the bash shell.
export SHELLOPTS := errexit

# Project main package location (can be multiple ones).
CMD_DIR := ./cmd

# Project output directory.
OUTPUT_DIR := ./bin

# Build direcotory.
BUILD_DIR := ./build

# Current version of the project.
VERSION ?= $(shell git describe --match 'v[0-9]*' --dirty='.m' --always)
GITSHA ?= $(shell git rev-parse --short HEAD)

# Track code version with Docker Label.
DOCKER_LABELS ?= git-describe="$(shell date -u +v%Y%m%d)-$(shell git describe --tags --always --dirty)"

# Golang standard bin directory.
GOPATH ?= $(shell go env GOPATH)
GOROOT ?= $(shell go env GOROOT)
BIN_DIR := $(GOROOT)/bin
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
.PHONY: lint test build container push addlicense

build: build-local

# more info about `GOGC` env: https://github.com/golangci/golangci-lint#memory-usage-of-golangci-lint
lint: $(GOLANGCI_LINT)
	@$(GOLANGCI_LINT) run

$(GOLANGCI_LINT):
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(BIN_DIR) v1.23.6

build-local:
	@for target in $(TARGETS); do                                                      \
	  CGO_ENABLED=$(CGO_ENABLED) go build -v -o $(OUTPUT_DIR)/$${target}               \
	    -ldflags "-s -w -X $(ROOT)/pkg/version.Version=$(VERSION)"                     \
	    $(CMD_DIR)/$${target};                                                         \
	done

addlicense:
	addlicense -c "The MIDI Authors" **/*.go

.PHONY: clean
clean:
	@-rm -vrf ${OUTPUT_DIR}
