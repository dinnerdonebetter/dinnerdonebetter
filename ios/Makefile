# iOS Project Makefile
# Provides convenient commands to run tests and build the iOS project

# Project configuration
PROJECT := ios.xcodeproj
SCHEME := ios
# DESTINATION can be overridden via environment variable (e.g., in CI)
# Default to iPhone 16 Pro with OS 18.2 for local development
DESTINATION ?= "platform=iOS Simulator,name=iPhone 16 Pro,OS=18.2"
TEST_TARGET := iosTests
UI_TEST_TARGET := iosUITests

# xcodebuild configuration
XCODEBUILD := xcodebuild
BUILD_DIR := build
DERIVED_DATA := $(BUILD_DIR)/DerivedData

.PHONY: test
test: unit-test # ui-test

.PHONY: unit-test
unit-test:
	$(XCODEBUILD) test \
		-project $(PROJECT) \
		-scheme $(SCHEME) \
		-destination $(DESTINATION) \
		-only-testing:$(TEST_TARGET) \
		-derivedDataPath $(DERIVED_DATA) \
		-quiet

.PHONY: ui-test
ui-test:
	$(XCODEBUILD) test \
		-project $(PROJECT) \
		-scheme $(SCHEME) \
		-destination $(DESTINATION) \
		-only-testing:$(UI_TEST_TARGET) \
		-derivedDataPath $(DERIVED_DATA) \
		-quiet

.PHONY: build
build:
	$(XCODEBUILD) build \
		-project $(PROJECT) \
		-scheme $(SCHEME) \
		-destination $(DESTINATION) \
		-derivedDataPath $(DERIVED_DATA) \
		-quiet

.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)
	$(XCODEBUILD) clean \
		-project $(PROJECT) \
		-scheme $(SCHEME) \
		-derivedDataPath $(DERIVED_DATA) \
		-quiet

.PHONY: lint
lint:
	@if ! command -v swiftlint >/dev/null 2>&1; then \
		echo "⚠️  SwiftLint not found. Installing via Homebrew..."; \
		brew install swiftlint || echo "❌ Failed to install SwiftLint. Please install manually: brew install swiftlint"; \
	fi
	swiftlint lint --strict

.PHONY: lint-autocorrect
lint-autocorrect:
	@if ! command -v swiftlint >/dev/null 2>&1; then \
		echo "⚠️  SwiftLint not found. Installing via Homebrew..."; \
		brew install swiftlint || echo "❌ Failed to install SwiftLint. Please install manually: brew install swiftlint"; \
	fi
	swiftlint lint --fix

.PHONY: list-destinations
list-destinations: ## List available iOS Simulator destinations
	xcrun simctl list devices available | grep -E "iPhone|iPad" || echo "No simulators found. Run 'xcrun simctl list devices' for more details."
