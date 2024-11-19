TEST_FOLDERS = ./pkg
EXAMPLE_FOLDER = ./cmd/examples
.PHONY: test

# run tests inside the TEST_FOLDERS
test:
	@echo Running tests in $(TEST_FOLDERS)... && \
	go test $(TEST_FOLDERS)/... -v || exit 1

# run example instance inside cmd/examples
run:
	@echo Running example inside $(EXAMPLE_FOLDER)
	go run $(EXAMPLE_FOLDER)|| exit 1
