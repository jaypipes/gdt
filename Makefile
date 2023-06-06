VERSION ?= $(shell git describe --tags --always --dirty)

.PHONY: test test-unit test-http

test-unit:
	go test -v ./
	go test -v ./http

test-http:
	go test -v ./examples/http/tests/api

test: test-unit test-http
