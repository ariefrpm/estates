

.PHONY: clean all init generate generate_mocks

all: build/main

build/main: cmd/main.go generated
	@echo "Building..."
	go build -o $@ $<

clean:
	rm -rf generated

init: clean generate
	go mod tidy
	go mod vendor

test:
	go clean -testcache
	go test -short -coverprofile coverage.out -short -v ./...
	grep -v -E "mock|api\.gen" coverage.out > coverage_filtered.out
	go tool cover -func=coverage_filtered.out

test_api:
	go clean -testcache
	go test ./tests/...

generate: generated generate_mocks

generated: api.yml
	@echo "Generating files..."
	mkdir generated || true
	oapi-codegen --package generated -generate types,server,spec $< > generated/api.gen.go

INTERFACES_GO_FILES := $(shell find core/interfaces -name "*_interface.go")
INTERFACES_GEN_GO_FILES := $(INTERFACES_GO_FILES:%.go=%_mock.go)

generate_mocks: $(INTERFACES_GEN_GO_FILES)

$(INTERFACES_GEN_GO_FILES): %_mock.go: %.go
	@echo "Generating mocks $@ for $<"
	mockgen -source=$< -destination=$@ -package=$(shell basename $(dir $<))
