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
DI_TARGETS             := internal/build/services/api \
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
LINTER_IMAGE           := golangci/golangci-lint:v1.64.5

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
	@mkdir --parents $(ARTIFACTS_DIR)

clean-$(ARTIFACTS_DIR):
	@rm -rf $(ARTIFACTS_DIR)

## PREREQUISITES

.PHONY: ensure_scc_installed
ensure_scc_installed:
ifeq (, $(shell which scc))
	$(shell go install github.com/boyter/scc@v3.4.0)
endif

.PHONY: ensure_wire_installed
ensure_wire_installed:
ifeq (, $(shell which wire))
	$(shell go install github.com/google/wire/cmd/wire@v0.6.0)
endif

.PHONY: ensure_fieldalignment_installed
ensure_fieldalignment_installed:
ifeq (, $(shell which fieldalignment))
	$(shell go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@v0.29.0)
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

.PHONY: ensure_sqlc_installed
ensure_sqlc_installed:
ifeq (, $(shell which sqlc))
	$(shell go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0)
endif

.PHONY: ensure_air_installed
ensure_air_installed:
ifeq (, $(shell which air))
	$(shell go install github.com/air-verse/air@v1.61.5)
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

.PHONY: setup
setup: $(ARTIFACTS_DIR) revendor rewire configs ensure_scc_installed ensure_wire_installed ensure_fieldalignment_installed ensure_tagalign_installed ensure_gci_installed ensure_sqlc_installed ensure_scc_installed

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
	for tgt in $(DI_TARGETS); do \
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
format_golang: format_imports format_go_fieldalignment # format_go_tag_alignment
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
golang_lint: check_lib_for_types sqlc_struct_check
	@$(CONTAINER_RUNNER) pull --quiet $(LINTER_IMAGE)
	$(LINTER) run --config=.golangci.yml --timeout 15m ./...

.PHONY: lint_terraform
lint_terraform: terraformat
	@(cd environments/dev/terraform && terraform init -upgrade && terraform validate && terraform fmt && terraform fmt -check)

.PHONY: lint
lint: lint_containers queries_lint golang_lint

.PHONY: clean_coverage
clean_coverage:
	@rm --force $(COVERAGE_OUT) profile.out;

.PHONY: coverage
coverage: clean_coverage $(ARTIFACTS_DIR)
	@$(GO_TEST) -coverprofile=$(COVERAGE_OUT) -covermode=atomic $(TESTABLE_PACKAGE_LIST) > /dev/null
	@go tool cover -func=$(ARTIFACTS_DIR)/coverage.out | grep 'total:' | xargs | awk '{ print "COVERAGE: " $$3 }'

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

.PHONY: codegen-ts
codegen-ts:
	go run $(THIS)/cmd/tools/codegen/openapi --typescript
	@(cd ../frontend && $(MAKE) format)

.PHONY: codegen-go
codegen-go:
	go run $(THIS)/cmd/tools/codegen/openapi --golang
	@$(MAKE) format_golang

.PHONY: codegen
codegen: codegen-ts codegen-go

.PHONY: querier
querier: queries queries_lint
	rm --recursive --force internal/database/postgres/generated/*.go
	$(SQL_GENERATOR) generate --no-remote

.PHONY: sqlc_struct_check
sqlc_struct_check:
	go run $(THIS)/cmd/tools/sqlc_struct_checker

.PHONY: openapi
openapi:
	go run $(THIS)/cmd/tools/codegen/openapi_spec
	(cd .. && $(MAKE) format_yaml)

.PHONY: env_vars
env_vars:
	go run $(THIS)/cmd/tools/codegen/valid_env_vars
	gofmt -s -w internal/config/envvars/*.go

.PHONY: proto
proto: ensure_protoc-gen-go_installed ensure_protoc-gen-go-grpc_installed
	rm -rf $(PROTO_OUTPUT_PATH)
	mkdir -p $(PROTO_OUTPUT_PATH)
	protoc --go_out=. \
		--go-grpc_out=. \
		--go_opt=module=$(THIS) \
		--go-grpc_opt=module=$(THIS) \
		-I internal/services/ \
		$(PROTO_FILES_PATH)
	$(MAKE) format_golang

## INTEGRATION TESTING

.PHONY: integration-tests
integration-tests: integration_tests

.PHONY: integration_tests
integration_tests: integration_tests_postgres

.PHONY: integration_tests_postgres
integration_tests_postgres:
	$(CONTAINER_RUNNER) compose \
	--file $(TEST_COMPOSE_FILES_DIR)/integration-tests.yaml build
	$(CONTAINER_RUNNER) compose \
	--file $(TEST_COMPOSE_FILES_DIR)/integration-tests.yaml up \
	--build \
	--quiet-pull \
	--force-recreate \
	--remove-orphans \
	--renew-anon-volumes \
	--always-recreate-deps \
	$(if $(filter y Y yes YES true TRUE plz sure yup YUP,$(LET_HANG)),,--exit-code-from tests) \
	--attach api_server \
	--attach tests

## EXECUTION

.PHONY: build
build:
	go build $(TOTAL_PACKAGE_LIST)

.PHONY: test
test: $(ARTIFACTS_DIR) vendor build
	$(GO_TEST) -failfast $(TESTABLE_PACKAGE_LIST)

.PHONY: dev
dev: $(ARTIFACTS_DIR)
	$(CONTAINER_RUNNER) compose \
	--file $(ENVIRONMENTS_DIR)/localdev/docker-compose.yaml build
	$(CONTAINER_RUNNER) compose \
	--file $(ENVIRONMENTS_DIR)/localdev/docker-compose.yaml up \
	--no-attach loki \
	--no-attach redis \
	--no-attach jaeger \
	--no-attach grafana \
	--no-attach postgres \
	--no-attach prometheus \
	--build \
	--quiet-pull \
	--force-recreate \
	--remove-orphans \
	--renew-anon-volumes \
	--always-recreate-deps \
	--abort-on-container-exit

.PHONY: run_admin_app
run_admin_app: ensure_air_installed
	(cd cmd/services/admin_webapp && CONFIGURATION_FILEPATH="../../../$(ENVIRONMENTS_DIR)/localdev/config_files/admin_webapp_config.json" air -c .air.toml)

## MISC

.PHONY: line_count
line_count: ensure_scc_installed
	@scc --include-ext go --exclude-dir vendor

# OPERATIONS

GCP_DEV_PROJECT_NAME := dinner-done-better-dev

# https://cloud.google.com/sql/docs/postgres/connect-admin-proxy#connect-tcp
.PHONY: proxy_dev_db
proxy_dev_db:
	cloud_sql_proxy $(GCP_DEV_PROJECT_NAME):us-central1:dev --port 5434 --gcloud-auth

.PHONY: deploy_dev
deploy_dev:
	skaffold run --filename=skaffold.yaml --build-concurrency 3 --profile dev
