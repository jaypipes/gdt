VERSION ?= $(shell git describe --tags --always --dirty)

.PHONY: test test-unit test-example ensure-mockery

generate-mocks: ${GOPATH}/bin/mockery
	GO111MODULE=on mockery -all -case underscore -keeptree -note "DO NOT EDIT MANUALLY. If you make changes to anything in interfaces.go, run make generate-mocks."

${GOPATH}/bin/mockery:
	GO111MODULE=on go get github.com/vektra/mockery/cmd/mockery

test-unit:
	GO111MODULE=on go test -v ./
	GO111MODULE=on go test -v ./http

test-example:
	GO111MODULE=on go test -v ./examples/books/api
	GO111MODULE=on go test -v ./examples/books/tests/api

test: test-unit test-example
