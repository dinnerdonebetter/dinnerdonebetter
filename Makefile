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

.PHONY: ensure_protoc_installed
ensure_protoc_installed:
ifeq (, $(shell which protoc-gen-go-grpc))
	$(shell brew install protobuf)
endif

.PHONY: ensure_protoc-gen-go_installed
ensure_protoc-gen-go_installed: ensure_protoc_installed
ifeq (, $(shell which protoc-gen-go-grpc))
	$(shell go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.4)
endif

.PHONY: ensure_protoc-gen-go-grpc_installed
ensure_protoc-gen-go-grpc_installed: ensure_protoc_installed
ifeq (, $(shell which protoc-gen-go-grpc))
	$(shell go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1)
endif

.PHONY: ensure_protoc-gen-swift_installed
ensure_protoc-gen-swift_installed: ensure_protoc_installed
ifeq (, $(shell which protoc-gen-swift))
	$(shell brew install swift-protobuf)
endif

.PHONY: ensure_protoc-gen-grpc-swift_installed
ensure_protoc-gen-grpc-swift_installed: ensure_protoc_installed
ifeq (, $(shell which protoc-gen-grpc-swift-2))
	$(shell brew install protoc-gen-grpc-swift)
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

.PHONY: regit
regit:
	cd ../
	git clone git@github.com:dinnerdonebetter/dinnerdonebetter tempdir
	@if [ -n "$(BRANCH)" ]; then \
	  (cd tempdir && git checkout $(BRANCH)); \
	fi
	cp -rf tempdir/.git .
	rm -rf tempdir

.PHONY: deploy_dev
deploy_dev:
	kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.16.2/cert-manager.yaml
	skaffold run --filename=skaffold.yaml --build-concurrency 1 --profile $(DEV_NAMESPACE)

.PHONY: nuke_dev
nuke_dev:
	kubectl delete deployments,cronjobs,configmaps,services,secrets --namespace $(DEV_NAMESPACE) --selector='managed_by!=terraform'

.PHONY: format_proto
format_proto:
	$(FORMAT_PROTOBUFS) format --path proto --write

# PATHS
PROTO_FILES_PATH          := proto/*.proto
PROTO_GO_OUTPUT_PATH      := backend
PROTO_OUTPUT_BACKEND_PATH := backend/internal/grpc
PROTO_OUTPUT_IOS_PATH     := ios/ios/Generated
BACKEND_REPO_NAME         := github.com/dinnerdonebetter/backend

.PHONY: proto_golang
proto_golang: ensure_protoc_installed ensure_protoc-gen-go_installed ensure_protoc-gen-go-grpc_installed
	mkdir -p $(PROTO_OUTPUT_BACKEND_PATH)
	protoc --go_out=$(PROTO_GO_OUTPUT_PATH) \
		--go-grpc_out=$(PROTO_GO_OUTPUT_PATH) \
		--go_opt=module=$(BACKEND_REPO_NAME) \
		--go-grpc_opt=module=$(BACKEND_REPO_NAME) \
		--proto_path proto/ \
		$(PROTO_FILES_PATH);
	(cd backend && $(MAKE) format_golang)

.PHONY: proto_swift
proto_swift: ensure_protoc-gen-swift_installed ensure_protoc-gen-grpc-swift_installed
	mkdir -p $(PROTO_OUTPUT_IOS_PATH)
	protoc --swift_out=$(PROTO_OUTPUT_IOS_PATH) \
		--grpc-swift-2_out=$(PROTO_OUTPUT_IOS_PATH) \
		--grpc-swift-2_opt=Client=true \
		--proto_path proto/ \
		$(PROTO_FILES_PATH)

.PHONY: proto
proto: format_proto proto_golang proto_swift
