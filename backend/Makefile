# ENVIRONMENT
PWD      := $(shell pwd)
MYSELF   := $(shell id -u)
MY_GROUP := $(shell id -g)

# PATHS
THIS                   := github.com/dinnerdonebetter/backend
ARTIFACTS_DIR          := artifacts
COVERAGE_OUT           := $(ARTIFACTS_DIR)/coverage.out
ENVIRONMENTS_DIR       := deploy/environments
TEST_COMPOSE_FILES_DIR := $(ENVIRONMENTS_DIR)/testing/compose_files
GENERATED_QUERIES_DIR  := internal/database/postgres/generated
PROTO_FILES_PATH       := proto/*.proto
PROTO_OUTPUT_PATH      := internal/grpc
WIRE_TARGETS           := internal/build/services/api/grpc \
	 internal/build/services/api/http \
	 internal/build/functions/data_change_message_handler \
	 internal/build/jobs/db_cleaner \
	 internal/build/jobs/search_data_index_scheduler \
	 internal/build/jobs/meal_plan_finalizer \
	 internal/build/jobs/meal_plan_grocery_list_initializer \
	 internal/build/jobs/meal_plan_task_creator

# COMPUTED
TOTAL_PACKAGE_LIST    := `go list $(THIS)/...`
TESTABLE_PACKAGE_LIST := `go list $(THIS)/... | grep -Ev '(cmd|integration|mock|fakes|converters|utils|generated)'`

# CONTAINER VERSIONS
CONTAINER_LINTER_IMAGE := openpolicyagent/conftest:v0.54.0
SQL_GENERATOR_IMAGE    := sqlc/sqlc:1.26.0
LINTER_IMAGE           := golangci/golangci-lint:v2.4.0

# COMMANDS
GO_FORMAT             := gofmt -s -w
GO_TEST               := CGO_ENABLED=1 go test -shuffle=on -race -vet=all
CONTAINER_RUNNER      := docker
RUN_CONTAINER         := $(CONTAINER_RUNNER) run --rm --volume $(PWD):$(PWD) --workdir=$(PWD)
RUN_CONTAINER_AS_USER := $(CONTAINER_RUNNER) run --rm --volume $(PWD):$(PWD) --workdir=$(PWD) --user $(MYSELF):$(MY_GROUP)
LINTER                := $(RUN_CONTAINER) $(LINTER_IMAGE) golangci-lint
SQL_GENERATOR         := $(RUN_CONTAINER_AS_USER) $(SQL_GENERATOR_IMAGE)

## non-PHONY folders/files

clean:
	rm -rf $(ARTIFACTS_DIR)

$(ARTIFACTS_DIR):
	@mkdir -p $(ARTIFACTS_DIR)

clean-$(ARTIFACTS_DIR):
	@rm -rf $(ARTIFACTS_DIR)

## PREREQUISITES

.PHONY: ensure_wire_installed
ensure_wire_installed:
ifeq (, $(shell which wire))
	$(shell go install github.com/google/wire/cmd/wire@v0.7.0)
endif

.PHONY: ensure_fieldalignment_installed
ensure_fieldalignment_installed:
ifeq (, $(shell which fieldalignment))
	$(shell go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@v0.37.0)
endif

.PHONY: ensure_tagalign_installed
ensure_tagalign_installed:
ifeq (, $(shell which tagalign))
	$(shell go install github.com/4meepo/tagalign/cmd/tagalign@v1.4.1)
endif

.PHONY: ensure_gci_installed
ensure_gci_installed:
ifeq (, $(shell which gci))
	$(shell go install github.com/daixiang0/gci@v0.13.5)
endif

.PHONY: ensure_goimports_installed
ensure_goimports_installed:
ifeq (, $(shell which goimports))
	$(shell go install golang.org/x/tools/cmd/goimports@v0.29.0)
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

.PHONY: ensure_air_installed
ensure_air_installed:
ifeq (, $(shell which air))
	$(shell go install github.com/air-verse/air@v1.63.0)
endif

.PHONY: setup
setup: $(ARTIFACTS_DIR) revendor rewire configs ensure_wire_installed ensure_fieldalignment_installed ensure_tagalign_installed ensure_gci_installed

.PHONY: clean_vendor
clean_vendor:
	rm -rf vendor go.sum

vendor:
	if [ ! -f go.mod ]; then go mod init; fi
	go mod tidy
	go mod vendor

.PHONY: revendor
revendor: clean_vendor vendor

## DEPENDENCY INJECTION

.PHONY: rewire
rewire: ensure_wire_installed
	for tgt in $(WIRE_TARGETS); do \
		rm -f $(THIS)/$$tgt/wire_gen.go && wire gen $(THIS)/$$tgt; \
	done

## FORMATTING

.PHONY: format_imports
format_imports: ensure_gci_installed
	gci write --section standard --section "prefix($(THIS))" --section "prefix($(dir $(THIS)))" --section default --custom-order `find $(PWD) -type f -not -path '*/vendor/*' -name "*.go"`

.PHONY: format_go_fieldalignment
format_go_fieldalignment: ensure_fieldalignment_installed
	@until fieldalignment -fix ./...; do true; done > /dev/null

.PHONY: format_go_tag_alignment
format_go_tag_alignment: ensure_tagalign_installed
	@until tagalign -fix -sort -order "env,envDefault,envPrefix,json,mapstructure,toml,yaml" ./...; do true; done > /dev/null

.PHONY: format_golang
format_golang: goimports format_imports format_go_fieldalignment # format_go_tag_alignment
	for file in `find $(PWD) -type f -not -path '*/vendor/*' -name "*.go"`; do $(GO_FORMAT) $$file; done

.PHONY: terraformat
terraformat:
	@(cd deploy/environments/dev/terraform && terraform fmt)

.PHONY: format
format: format_golang terraformat

.PHONY: fmt
fmt: format

.PHONY: goimports
goimports: ensure_goimports_installed
	goimports -w .

## LINTING

.PHONY: lint_containers
lint_containers:
	@$(CONTAINER_RUNNER) pull --quiet $(CONTAINER_LINTER_IMAGE)
	$(RUN_CONTAINER_AS_USER) $(CONTAINER_LINTER_IMAGE) test --policy containers.rego `find . -type f -name "*.Dockerfile"`

.PHONY: queries_lint
queries_lint:
	@$(CONTAINER_RUNNER) pull --quiet $(SQL_GENERATOR_IMAGE)
	$(SQL_GENERATOR) compile --no-remote
	$(SQL_GENERATOR) vet --no-remote

.PHONY: golang_lint
golang_lint: struct_usage_check
	@$(CONTAINER_RUNNER) pull --quiet $(LINTER_IMAGE)
	$(LINTER) run --config=.golangci.yml --timeout 30m ./...

.PHONY: lint_terraform
lint_terraform: terraformat
	@(cd environments/dev/terraform && terraform init -upgrade && terraform validate && terraform fmt && terraform fmt -check)

.PHONY: lint
lint: lint_containers queries_lint golang_lint

.PHONY: check_lib_for_types
check_lib_for_types:
	@! grep -r 'pkg/types"' internal/lib/
	@! grep -r 'pkg/types/fakes"' internal/lib/

## GENERATED FILES

.PHONY: configs
configs:
	go run $(THIS)/cmd/tools/codegen/configs

.PHONY: queries
queries:
	go run $(THIS)/cmd/tools/codegen/queries

.PHONY: querier
querier: queries queries_lint
	rm -rf internal/database/postgres/generated/*.go
	$(SQL_GENERATOR) generate --no-remote

.PHONY: struct_usage_check
struct_usage_check:
	go run $(THIS)/cmd/tools/struct_usage_checker

.PHONY: env_vars
env_vars:
	go run $(THIS)/cmd/tools/codegen/valid_env_vars
	gofmt -s -w internal/config/envvars/*.go

## INTEGRATION TESTING

.PHONY: integration-tests
integration-tests: integration_tests

.PHONY: integration_tests
integration_tests: integration_tests_postgres

.PHONY: integration_tests_postgres
integration_tests_postgres:
	go test -v -count=1 $(THIS)/tests_integration/apiserver

## EXECUTION

.PHONY: build
build:
	go build $(TOTAL_PACKAGE_LIST)

.PHONY: test
test: $(ARTIFACTS_DIR) vendor build
	$(GO_TEST) -failfast $(TESTABLE_PACKAGE_LIST)

.PHONY: dev
dev: $(ARTIFACTS_DIR)
	go run $(THIS)/cmd/localdev/server

AIR_PORT := 9999

.PHONY: admin
admin: $(ARTIFACTS_DIR) ensure_air_installed
	air --build.cmd "go build -o ./$(ARTIFACTS_DIR)/admin_tmp $(THIS)/cmd/services/admin" --build.bin "./$(ARTIFACTS_DIR)/admin_tmp" --proxy.app_port "8888" -proxy.proxy_port "9999" -proxy.enabled "true" --build.log "./$(ARTIFACTS_DIR)/build-errors.log"

# OPERATIONS

GCP_DEV_PROJECT_NAME := dinner-done-better-dev

# https://cloud.google.com/sql/docs/postgres/connect-admin-proxy#connect-tcp
.PHONY: proxy_dev_db
proxy_dev_db:
	cloud_sql_proxy $(GCP_DEV_PROJECT_NAME):us-central1:dev --port 5434 --gcloud-auth

.PHONY: deploy_dev
deploy_dev:
	skaffold run --filename=skaffold.yaml --build-concurrency 3 --profile dev
