.SILENT = all clean test install
.PHONY: build

ARCHS = arm64 amd64 386 arm
PLATFORMS = linux

PREFIX = /usr/local
PLUGIN_DIR = $(PREFIX)/lib/system-confd/plugins

all: build

build:
	@echo "Building the plugin..."
	@# Add build commands here
	@for arch in $(ARCHS); do \
		for platform in $(PLATFORMS); do \
			echo "Building for $$platform on $$arch..."; \
			# Simulate build command \
			mkdir -p build/$$platform-$$arch; \
			GOARCH=$$arch GOOS=$$platform go build -o build/$$platform-$$arch/keyboard.plugin; \
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

install:
	@echo "Installing keyboard plugin..."
	@# Add install commands here
	@install -d -v -m 755 -o root -g root $(PLUGIN_DIR)
	@if [[ "$(shell uname -m)" == "x86_64" ]]; then \
	    ARCH="amd64"; \
	elif [[ "$(shell uname -m)" == "aarch64" ]]; then \
	    ARCH="arm64"; \
	elif [[ "$(shell uname -m)" == "i386" || "$(shell uname -m)" == "i586" || "$(shell uname -m)" == "i686" ]]; then \
	    ARCH="386"; \
	elif [[ "$(shell uname -m)" == "armv7l" ]]; then \
	    ARCH="arm"; \
	else \
	    echo "Unsupported architecture: $(shell uname -m)"; exit 1; \
	fi; \
	@install -v -m 755 -o root -g root build/linux-$(ARCH)/keyboard.plugin $(PLUGIN_DIR)/keyboard.plugin
	@echo "Installation completed."
