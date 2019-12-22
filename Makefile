VERSION ?= $(shell git describe --tags --always --dirty)

.PHONY: test test-unit test-example

test-unit:
	GO111MODULE=on go test -v ./
	GO111MODULE=on go test -v ./http

test-example:
	GO111MODULE=on go test -v ./examples/books/api
	GO111MODULE=on go test -v ./examples/books/tests/api

test: test-unit test-example
