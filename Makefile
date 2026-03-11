PWD           := $(shell pwd)
MYSELF        := $(shell id -u)
MY_GROUP      := $(shell id -g)

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

.PHONY: lint_markdown
lint_markdown:
	./scripts/lint_markdown.sh

.PHONY: setup
setup: ensure_yamlfmt_installed
	(cd backend && $(MAKE) setup)

.PHONY: format
format: format_yaml
	(cd backend && $(MAKE) format)
	(cd ios && $(MAKE) format)

.PHONY: terraformat
terraformat:
	(cd backend && $(MAKE) terraformat)

.PHONY: format_yaml
format_yaml: ensure_yamlfmt_installed
	yamlfmt -conf .yamlfmt.yaml

.PHONY: lint
lint:
	(cd backend && $(MAKE) lint)
	(cd ios && $(MAKE) lint)

.PHONY: test
test:
	(cd backend && $(MAKE) test)
	(cd ios && $(MAKE) test)

.PHONY: regit
regit:
	cd ../
	git clone git@github.com:dinnerdonebetter/dinnerdonebetter tempdir
	@if [ -n "$(BRANCH)" ]; then \
	  (cd tempdir && git checkout $(BRANCH)); \
	fi
	cp -rf tempdir/.git .
	rm -rf tempdir

# Deploy to Docker Desktop Kubernetes cluster (no Helm required)
.PHONY: deploy_localdev
deploy_localdev:
	@echo "Deploying infrastructure to Docker Desktop Kubernetes..."
	skaffold run --filename=infra/skaffold.yaml --build-concurrency 1 --profile localdev
	@echo "Deploying backend services to Docker Desktop Kubernetes..."
	KO_DOCKER_REPO=ko.local skaffold run --filename=backend/skaffold.yaml --build-concurrency 1 --profile localdev
	@echo "Deployment complete! Services available at:"
	@echo "  - API Server: http://localhost:8000"
	@echo "  - Admin Webapp: http://localhost:8888"

# Deploy prod Terraform: infra first (GKE, networking), then backend. Run from repo root.
# Pass args through, e.g. make deploy_terraform_prod ARGS="-auto-approve"
.PHONY: deploy_terraform_prod
deploy_terraform_prod:
	./infra/scripts/terraform_apply_prod.sh $(ARGS)
	./backend/scripts/terraform_apply_prod.sh $(ARGS)

# Prod deploy + verify. Run from repo root. Requires kubectl pointed at prod, grpcurl.
.PHONY: deploy_prod
deploy_prod:
	./scripts/deploy-prod-local.sh

.PHONY: verify_prod
verify_prod:
	skaffold verify --filename=skaffold.yaml --profile prod
	
.PHONY: deploy_prod_verify
deploy_prod_verify: deploy_prod
	@echo "Waiting 60s for load balancer..."
	@sleep 60
	$(MAKE) verify_prod

.PHONY: nuke_localdev
nuke_localdev:
	kubectl delete namespace localdev --ignore-not-found

.PHONY: format_proto
format_proto:
	$(FORMAT_PROTOBUFS) format --path proto --write

# PATHS
# Exclude monolithic {domain}/{domain}.proto files (they're duplicates of the split files)
# But keep common.proto and filtering.proto at the root, and uploaded_media.proto (hasn't been split yet)
PROTO_FILES_PATH          := $(shell find proto -name "*.proto" -type f ! -path "proto/audit/audit.proto" ! -path "proto/auth/auth.proto" ! -path "proto/dataprivacy/dataprivacy.proto" ! -path "proto/identity/identity.proto" ! -path "proto/internal_ops/internal_ops.proto" ! -path "proto/issue_reports/issue_reports.proto" ! -path "proto/mealplanning/mealplanning.proto" ! -path "proto/notifications/notifications.proto" ! -path "proto/oauth/oauth.proto" ! -path "proto/settings/settings.proto" ! -path "proto/waitlists/waitlists.proto" ! -path "proto/webhooks/webhooks.proto")
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
      	--grpc-swift-2_opt=Client=true,Server=false \
      	--swift_opt=Visibility=Public \
		--proto_path proto/ \
		$(PROTO_FILES_PATH)
	(cd ios && $(MAKE) format)

.PHONY: proto
proto: format_proto proto_golang proto_swift
