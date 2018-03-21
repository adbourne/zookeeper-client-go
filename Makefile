# Determine this makefile's path.
# Be sure to place this BEFORE `include` directives, if any.
THIS_FILE := $(lastword $(MAKEFILE_LIST))

GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)

default: dependencies buildTools lint build test

clean:
	rm -rf zookeeper/bin/*; rm -rf pkg/*; rm -rf vendor/*

## Update/download dependencies
dependencies:
	dep ensure

## Build tools ensures that the tools used in the build toolchain are installed and configured
## this should only have to be run once
buildTools:
	go get -u github.com/kyoh86/richgo && \
	go get -u github.com/alecthomas/gometalinter \
	&& gometalinter --install

## Runs go fmt on the entire project, excluding the vendor directory
fmt:
	gofmt -w $(GOFMT_FILES)

## Runs the gometalinter on the entire project, excluding the vendor directory
lint:
	gometalinter ./... | grep -v vendor/ | sed ''/warning/s//$$(printf "\033[33mwarning\033[0m")/'' | sed ''/error/s//$$(printf "\033[31merror\033[0m")/''

## Builds the project
build:
	go install ./zookeeper

## Generates mocks used for testing
mocks:
	mockery -dir services/ -all -case underscore

## Runs all tests in the project, excluding the vendor directory
test:
	richgo test ./... -v --cover

.PHONY: test package