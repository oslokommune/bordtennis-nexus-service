GOPATH := $(shell go env GOPATH)
GOBIN  ?= $(GOPATH)/bin

RICHGO := $(GOBIN)/richgo

$(RICHGO):
	go install github.com/kyoh86/richgo@v0.3.11

test: $(RICHGO)
	$(RICHGO) test ./...

check:
	npx pre-commit run

run:
	go run *.go
