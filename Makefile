setup:
	go get -u github.com/golang/lint/golint

lint:
	golint -set_exit_status $(shell go list ./... | grep -v /vendor/)

test:
	go test $(shell go list ./... | grep -v /vendor/)

vet:
	go vet $(shell go list ./... | grep -v /vendor/)

test-all: vet lint test

.PHONY: setup lint test vet test-all
