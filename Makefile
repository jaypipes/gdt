VERSION ?= $(shell git describe --tags --always --dirty)

.PHONY: test test-unit test-example ensure-mockery

test-unit:
	go test -v ./
	go test -v ./http

test-example:
	go test -v ./examples/books/api
	go test -v ./examples/books/tests/api

test: test-unit test-example
