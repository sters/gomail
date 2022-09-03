export GOBIN := $(PWD)/bin
export PATH := $(GOBIN):$(PATH)

TOOLS=$(shell cat tools/tools.go | egrep '^\s_ '  | awk '{ print $$2 }')

.PHONY: bootstrap-tools
bootstrap-tools:
	@echo "Installing: " $(TOOLS)
	@go install $(TOOLS)

.PHONY: lint
lint: bootstrap-tools
	$(GOBIN)/golangci-lint run -v ./...

.PHONY: lint-fix
lint-fix: bootstrap-tools
	$(GOBIN)/golangci-lint run --fix -v ./...

.PHONY: test
test:
	go test -v -race ./...

.PHONY: coverage
coverage: coverage.txt

coverage.txt: $(GOFILES) Makefile
	go test -coverprofile=coverage.txt -covermode=atomic ./...

.PHONY: show-coverage
show-coverage: coverage.txt
	go tool cover -html=coverage.txt

.PHONY: clean
clean:
	go clean
	rm -f coverage.txt

.PHONY: tidy
tidy:
	go mod tidy