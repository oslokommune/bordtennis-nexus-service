GOPATH := $(shell go env GOPATH)
GOBIN  ?= $(GOPATH)/bin

RICHGO := $(GOBIN)/richgo
GOLANGCI_LINT := $(GOBIN)/golangci-lint
GOFUMPT := $(GOBIN)/gofumpt

$(RICHGO):
	go install github.com/kyoh86/richgo@v0.3.11
$(GOLANGCI_LINT):
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.1
$(GOFUMPT):
	go install mvdan.cc/gofumpt@0.4.0

test: $(RICHGO)
	$(RICHGO) test ./...

lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run

fmt: $(GOFUMPT)
	$(GOFUMPT) -w .

ci: fmt
	git diff --exit-code > /dev/null

check: fmt lint test

run:
	go run *.go
