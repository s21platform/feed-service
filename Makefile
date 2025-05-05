.PHONY: lint protogen

GOLANGCI_LINT_INSTALL_DIR ?= $(shell go env GOPATH)/bin
GOLANGCI_LINT := $(GOLANGCI_LINT_INSTALL_DIR)/golangci-lint

lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run

$(GOLANGCI_LINT):
	@echo "golangci-lint not found, installing the latest version..."
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOLANGCI_LINT_INSTALL_DIR)

protogen:
	protoc --go_out=. --go-grpc_out=. ./api/feed.proto
	protoc --doc_out=. --doc_opt=markdown,GRPC_API.md ./api/feed.proto

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out