.PHONY: validate-examples
validate-examples: ## validate examples in the specification markdown files
	go test ./schema/example_test.go

.PHONY: test
test:
	go test ./...

.PHONY: generate-python-models
generate-python-models: ## generate Python models from JSON schema
	python3 tools/generate_python_models.py
