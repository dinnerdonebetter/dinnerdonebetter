PWD           := $(shell pwd)
MYSELF        := $(shell id -u)
MY_GROUP      := $(shell id -g)
DEV_NAMESPACE := dev

# PATHS
PROTO_FILES_PATH       := proto/*.proto
PROTO_OUTPUT_BACKEND_PATH      := backend/internal/grpc
BACKEND_REPO_NAME                   := github.com/dinnerdonebetter/backend

.PHONY: ensure_yamlfmt_installed
ensure_yamlfmt_installed:
ifeq (, $(shell which yamlfmt))
	$(shell go install github.com/google/yamlfmt/cmd/yamlfmt@latest)
endif

.PHONY: ensure_protoc-gen-go_installed
ensure_protoc-gen-go_installed:
ifeq (, $(shell which protoc-gen-go-grpc))
	$(shell go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.4)
endif

.PHONY: ensure_protoc-gen-go-grpc_installed
ensure_protoc-gen-go-grpc_installed:
ifeq (, $(shell which protoc-gen-go-grpc))
	$(shell go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1)
endif

.PHONY: setup
setup: ensure_yamlfmt_installed
	(cd backend && $(MAKE) setup)
	(cd frontend && $(MAKE) setup)

.PHONY: format
format: format_yaml
	(cd backend && $(MAKE) format)
	(cd frontend && $(MAKE) format)

.PHONY: terraformat
terraformat:
	(cd backend && $(MAKE) terraformat)
	(cd frontend && $(MAKE) terraformat)
	(cd infra && $(MAKE) terraformat)

.PHONY: format_yaml
format_yaml: ensure_yamlfmt_installed
	yamlfmt -conf .yamlfmt.yaml

.PHONY: lint
lint:
	(cd backend && $(MAKE) lint)
	(cd frontend && $(MAKE) lint)

.PHONY: test
test:
	(cd backend && $(MAKE) test)
	(cd frontend && $(MAKE) test)

.PHONY: proto
proto: backend_proto

.PHONY: backend_proto
backend_proto: ensure_protoc-gen-go_installed ensure_protoc-gen-go-grpc_installed
	rm -rf $(PROTO_OUTPUT_BACKEND_PATH)
	mkdir -p $(PROTO_OUTPUT_BACKEND_PATH)
	protoc --go_out=. \
		--go-grpc_out=. \
		--go_opt=module=$(BACKEND_REPO_NAME) \
		--go-grpc_opt=module=$(BACKEND_REPO_NAME) \
		-I internal/services/ \
		$(PROTO_FILES_PATH)
