PWD           := $(shell pwd)
MYSELF        := $(shell id -u)
MY_GROUP      := $(shell id -g)
DEV_NAMESPACE := dev

# CONTAINER VERSIONS
PROTOBUF_FORMAT := bufbuild/buf:1.5.0

# COMMANDS
CONTAINER_RUNNER      := docker
RUN_CONTAINER         := $(CONTAINER_RUNNER) run --rm --volume $(PWD):$(PWD) --workdir=$(PWD)
RUN_CONTAINER_AS_USER := $(CONTAINER_RUNNER) run --rm --volume $(PWD):$(PWD) --workdir=$(PWD) --user $(MYSELF):$(MY_GROUP)
FORMAT_PROTOBUFS      := $(RUN_CONTAINER) $(PROTOBUF_FORMAT)

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

.PHONY: format_proto
format_proto:
	$(FORMAT_PROTOBUFS) format --path proto --write

# PATHS
PROTO_FILES_PATH          := proto/*.proto
PROTO_GO_OUTPUT_PATH      := backend
PROTO_OUTPUT_BACKEND_PATH := backend/internal/grpc
BACKEND_REPO_NAME         := github.com/dinnerdonebetter/backend
GRPC_SERVICES             := core eating

.PHONY: backend_proto
backend_proto: ensure_protoc-gen-go_installed ensure_protoc-gen-go-grpc_installed format_proto
	mkdir -p $(PROTO_OUTPUT_BACKEND_PATH)
	for svc in $(GRPC_SERVICES); do \
		protoc --go_out=$(PROTO_GO_OUTPUT_PATH) \
			--go-grpc_out=$(PROTO_GO_OUTPUT_PATH) \
			--go_opt=module=$(BACKEND_REPO_NAME) \
			--go-grpc_opt=module=$(BACKEND_REPO_NAME) \
			--proto_path proto/ \
			$(PROTO_FILES_PATH); \
	done
	(cd backend && $(MAKE) format)

.PHONY: proto
proto: backend_proto
