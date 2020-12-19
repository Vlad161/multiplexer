TIMEOUT  := 1s
HOST     := localhost:8081

.PHONY: deps
deps:
	@go mod tidy && go mod vendor && go mod verify

.PHONY: test
test:
	@go test -race -count 1 -timeout $(TIMEOUT) ./...
