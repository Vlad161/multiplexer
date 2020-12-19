HOST := localhost:8081

PWD := $(PWD)
export PATH := $(PWD)/bin:$(PATH)

.PHONY: deps
deps:
	@go mod tidy && go mod vendor && go mod verify

.PHONY: tools
tools:
	cd tools && go generate -tags tools

.PHONY: test
test:
	@go test -race -count 1 ./...

.PHONY: generate
generate:
	@go generate ./...
