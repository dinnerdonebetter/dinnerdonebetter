PWD                    := $(shell pwd)
MYSELF                 := $(shell id -u)
MY_GROUP               := $(shell id -g)
GOPATH                 := $(GOPATH)
ARTIFACTS_DIR          := artifacts
COVERAGE_OUT           := $(ARTIFACTS_DIR)/coverage.out
GO_FORMAT              := gofmt -s -w
THIS                   := github.com/dinnerdonebetter/backend
TOTAL_PACKAGE_LIST     := `go list $(THIS)/...`
TESTABLE_PACKAGE_LIST  := `go list $(THIS)/... | grep -Ev '(cmd|integration|mock|fakes|converters|utils|generated)'`
ENVIRONMENTS_DIR       := environments
TEST_COMPOSE_FILES_DIR := $(ENVIRONMENTS_DIR)/testing/compose_files
GENERATED_QUERIES_DIR  := internal/database/postgres/generated
CLOUD_JOBS             := meal_plan_finalizer \
						  meal_plan_grocery_list_initializer \
						  meal_plan_task_creator \
						  search_data_index_scheduler \
						  email_prober
CLOUD_FUNCTIONS        := data_changes \
						  outbound_emailer \
						  search_indexer \
						  webhook_executor \
						  user_data_aggregator
WIRE_TARGETS           := server/http/build
CONTAINER_RUNNER       := docker
SQL_GENERATOR_IMAGE    := sqlc/sqlc:1.26.0
LINTER_IMAGE           := golangci/golangci-lint:v1.61.0
LINTER                 := $(CONTAINER_RUNNER) run --rm \
						  --volume $(PWD):$(PWD) \
						  --workdir=$(PWD) \
						  $(LINTER_IMAGE) golangci-lint
CONTAINER_LINTER_IMAGE := openpolicyagent/conftest:v0.54.0
SQL_GENERATOR          := $(CONTAINER_RUNNER) run --rm \
						  --volume $(PWD):/src \
						  --workdir /src \
						  --user $(MYSELF):$(MY_GROUP) \
						  $(SQL_GENERATOR_IMAGE)

## non-PHONY folders/files

clean:
	rm -rf $(ARTIFACTS_DIR)

$(ARTIFACTS_DIR):
	@mkdir --parents $(ARTIFACTS_DIR)

clean-$(ARTIFACTS_DIR):
	@rm -rf $(ARTIFACTS_DIR)

## prerequisites

.PHONY: setup
setup: $(ARTIFACTS_DIR) revendor rewire configs

.PHONY: ensure_wire_installed
ensure_wire_installed:
ifeq (, $(shell which wire))
	$(shell go install github.com/google/wire/cmd/wire@latest)
endif

.PHONY: ensure_fieldalignment_installed
ensure_fieldalignment_installed:
ifeq (, $(shell which fieldalignment))
	$(shell go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest)
endif

.PHONY: ensure_tagalign_installed
ensure_tagalign_installed:
ifeq (, $(shell which tagalign))
	$(shell go install github.com/4meepo/tagalign/cmd/tagalign@latest)
endif

.PHONY: ensure_gci_installed
ensure_gci_installed:
ifeq (, $(shell which gci))
	$(shell go install github.com/daixiang0/gci@latest)
endif

.PHONY: ensure_sqlc_installed
ensure_sqlc_installed:
ifeq (, $(shell which sqlc))
	$(shell go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest)
endif

.PHONY: ensure_scc_installed
ensure_scc_installed:
ifeq (, $(shell which scc))
	$(shell go install github.com/boyter/scc@latest)
endif

.PHONY: clean_vendor
clean_vendor:
	rm -rf vendor go.sum

vendor:
	if [ ! -f go.mod ]; then go mod init; fi
	go mod tidy
	go mod vendor
	for thing in $(CLOUD_FUNCTIONS); do \
  		(cd cmd/functions/$$thing && go mod tidy) \
	done
	for thing in $(CLOUD_JOBS); do \
  		(cd cmd/jobs/$$thing && go mod tidy) \
	done

.PHONY: revendor
revendor: clean_vendor vendor

## dependency injection

.PHONY: rewire
rewire:
	for tgt in $(WIRE_TARGETS); do \
		rm -f $(THIS)/internal/$$tgt/wire_gen.go && wire gen $(THIS)/internal/$$tgt; \
	done

## formatting

.PHONY: format
format: format_golang terraformat

.PHONY: format_golang
format_golang: format_imports ensure_fieldalignment_installed ensure_tagalign_installed
	@until fieldalignment -fix ./...; do true; done > /dev/null
	@until tagalign -fix -sort -order "json,toml,yaml" ./...; do true; done > /dev/null
	for file in `find $(PWD) -type f -not -path '*/vendor/*' -name "*.go"`; do $(GO_FORMAT) $$file; done

.PHONY: format_imports
format_imports: ensure_gci_installed
	gci write --section standard --section "prefix($(THIS))" --section "prefix($(dir $(THIS)))" --section default --custom-order `find $(PWD) -type f -not -path '*/vendor/*' -name "*.go"`

.PHONY: terraformat
terraformat:
	@touch environments/dev/terraform/service-config.json
	@touch environments/dev/terraform/worker-config.json
	@(cd environments/dev/terraform && terraform fmt)

.PHONY: lint_terraform
lint_terraform: terraformat
	@(cd environments/dev/terraform && terraform init -upgrade && terraform validate && terraform fmt && terraform fmt -check)

.PHONY: fmt
fmt: format terraformat

## Testing things

.PHONY: pre_lint
pre_lint:
	@until fieldalignment -fix ./...; do true; done > /dev/null
	@echo ""

.PHONY: lint_containers
lint_containers:
	@$(CONTAINER_RUNNER) pull --quiet $(CONTAINER_LINTER_IMAGE)
	$(CONTAINER_RUNNER) run --rm --volume $(PWD):$(PWD) --workdir=$(PWD) --user $(MYSELF):$(MY_GROUP) $(CONTAINER_LINTER_IMAGE) test --policy containers.rego `find . -type f -name "*.Dockerfile"`

.PHONY: queries_lint
queries_lint:
	@$(CONTAINER_RUNNER) pull --quiet $(SQL_GENERATOR_IMAGE)
	$(SQL_GENERATOR) compile --no-remote
	$(SQL_GENERATOR) vet --no-remote

.PHONY: golang_lint
golang_lint: sqlc_struct_check
	@$(CONTAINER_RUNNER) pull --quiet $(LINTER_IMAGE)
	$(LINTER) run --config=.golangci.yml --timeout 15m ./...

.PHONY: lint
lint: lint_containers queries_lint golang_lint

.PHONY: clean_coverage
clean_coverage:
	@rm --force $(COVERAGE_OUT) profile.out;

.PHONY: coverage
coverage: clean_coverage $(ARTIFACTS_DIR)
	@go test -coverprofile=$(COVERAGE_OUT) -shuffle=on -covermode=atomic -race $(TESTABLE_PACKAGE_LIST) > /dev/null
	@go tool cover -func=$(ARTIFACTS_DIR)/coverage.out | grep 'total:' | xargs | awk '{ print "COVERAGE: " $$3 }'

.PHONY: build
build:
	go build $(TOTAL_PACKAGE_LIST)

.PHONY: test
test: $(ARTIFACTS_DIR) vendor build
	go test -cover -shuffle=on -race -vet=all -failfast $(TESTABLE_PACKAGE_LIST)

## Generated files

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
querier: queries queries_lint  ## TODO: Why is this target broken?
	rm --recursive --force internal/database/postgres/generated/*.go
	$(SQL_GENERATOR) generate --no-remote

.PHONY: sqlc_struct_check
sqlc_struct_check:
	go run $(THIS)/cmd/tools/sqlc_struct_checker

.PHONY: openapi
openapi:
	go run $(THIS)/cmd/tools/codegen/openapi_spec

.PHONY: openapi-client
openapi-client: ensure_oapi_codegen_installed
	@sed -i 's/openapi: 3.1.0/openapi: 3.0.0/g' ../openapi_spec.yaml
	oapi-codegen --config=oapi-codegen.yaml ../openapi_spec.yaml
	@sed -i 's/openapi: 3.0.0/openapi: 3.1.0/g' ../openapi_spec.yaml

## Integration tests

.PHONY: wipe_containers
wipe_containers:
	@$(CONTAINER_RUNNER) stop $(shell $(CONTAINER_RUNNER) ps -aq) && $(CONTAINER_RUNNER) rm $(shell $(CONTAINER_RUNNER) ps -aq)

.PHONY: container_wipe
container_wipe: wipe_containers

.PHONY: integration-tests
integration-tests: integration_tests

.PHONY: integration_tests
integration_tests: integration_tests_postgres

.PHONY: integration_tests_postgres
integration_tests_postgres:
	$(CONTAINER_RUNNER) compose \
	--file $(TEST_COMPOSE_FILES_DIR)/integration-tests.yaml up \
	--build \
	--force-recreate \
	--remove-orphans \
	--renew-anon-volumes \
	--always-recreate-deps \
	$(if $(filter y Y yes YES true TRUE plz sure yup YUP,$(LET_HANG)),,--exit-code-from tests) \
	--attach api_server \
	--attach tests

## Running

.PHONY: dev
dev: $(ARTIFACTS_DIR)
	$(CONTAINER_RUNNER) compose \
	--file $(ENVIRONMENTS_DIR)/local/compose_files/dev.yaml up \
	--quiet-pull \
	--always-recreate-deps

## misc

.PHONY: tree
tree:
	# there are no long args for tree, but d means "directories only" and I means "ignore pattern"
	tree -d -I vendor

.PHONY: line_count
line_count: ensure_scc_installed
	@scc --include-ext go --exclude-dir vendor

## maintenance

# https://cloud.google.com/sql/docs/postgres/connect-admin-proxy#connect-tcp
.PHONY: proxy_dev_db
proxy_dev_db:
	cloud_sql_proxy dinner-done-better-dev:us-central1:dev --port 5434 --gcloud-auth
