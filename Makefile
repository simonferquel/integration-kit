DEV_IMAGE?=integration-kit-dev

setup:
	go get -u github.com/golang/lint/golint

lint:
	golint -set_exit_status $(shell go list ./... | grep -v /vendor/)

test:
	go test $(shell go list ./... | grep -v /vendor/)

vet:
	go vet $(shell go list ./... | grep -v /vendor/)

test-all: vet lint test

clean-dev-image:
	docker image rm ${DEV_IMAGE}

clean: clean-dev-image

.PHONY: setup lint test vet test-all clean-dev-image clean
