.SILENT = all clean test
.PHONY: build

ARCHS = arm64 amd64 386 arm
PLATFORMS = linux

all: build

build:
	@echo "Building the project..."
	@# Add build commands here
	@for arch in $(ARCHS); do \
		for platform in $(PLATFORMS); do \
			echo "Building for $$platform on $$arch..."; \
			# Simulate build command \
			mkdir -p build/$$platform-$$arch; \
			GOARCH=$$arch GOOS=$$platform go build -o build/$$platform-$$arch/systemconfd; \
		done; \
	done
	@echo "Build completed."

clean:
	@echo "Cleaning up..."
	@# Add clean commands here
	@rm -rfv build/
	@echo "Clean completed."

test:
	@echo "Running tests..."
	@# Add test commands here
	@go test ./...
	@gosec ./...
	@echo "Tests completed."

