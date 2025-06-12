.PHONY: validate-examples
validate-examples: ## validate examples in the specification markdown files
	go test ./schema/example_test.go

.PHONY: test
test:
	go test ./...
