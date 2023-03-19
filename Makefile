PWD                           := $(shell pwd)
GOPATH                        := $(GOPATH)
ARTIFACTS_DIR                 := artifacts
COVERAGE_OUT                  := $(ARTIFACTS_DIR)/coverage.out
GO                            := docker run --interactive --tty --volume $(PWD):$(PWD) --workdir $(PWD) --user $(shell id -u):$(shell id -g) golang:1.18-stretch go
GO_FORMAT                     := gofmt -s -w
THIS                          := github.com/prixfixeco/backend
TOTAL_PACKAGE_LIST            := `go list $(THIS)/...`
TESTABLE_PACKAGE_LIST         := `go list $(THIS)/... | grep -Ev '(cmd|tests|testutil|mock|fake)'`
ENVIRONMENTS_DIR              := environments
TEST_ENVIRONMENT_DIR          := $(ENVIRONMENTS_DIR)/testing
TEST_DOCKER_COMPOSE_FILES_DIR := $(TEST_ENVIRONMENT_DIR)/compose_files
LOCAL_ADDRESS                 := api.prixfixe.local
DEFAULT_CERT_TARGETS          := $(LOCAL_ADDRESS) prixfixe.local localhost 127.0.0.1 ::1
SQL_GENERATOR                 := docker run --rm --volume `pwd`:/src --workdir /src kjconroy/sqlc:1.16.0
GENERATED_QUERIES_DIR         := internal/database/postgres/generated
CLOUD_FUNCTIONS               := data_changes meal_plan_finalizer meal_plan_grocery_list_initializer meal_plan_task_creator

## non-PHONY folders/files

clear:
	@printf "\033[2J\033[3J\033[1;1H"

clean:
	rm -rf $(ARTIFACTS_DIR)

$(ARTIFACTS_DIR):
	@mkdir --parents $(ARTIFACTS_DIR)

clean-$(ARTIFACTS_DIR):
	@rm -rf $(ARTIFACTS_DIR)

.PHONY: setup
setup: $(ARTIFACTS_DIR) revendor rewire configs

## prerequisites

# not a bad idea to do this either:
## GO111MODULE=off go install golang.org/x/tools/...

ensure_wire_installed:
ifndef $(shell command -v wire 2> /dev/null)
	$(shell GO111MODULE=off go install github.com/google/wire/cmd/wire@latest)
endif

ensure_fieldalign_installed:
ifndef $(shell command -v wire 2> /dev/null)
	$(shell GO111MODULE=off go get -u golang.org/x/tools/...)
endif

ensure_scc_installed:
ifndef $(shell command -v scc 2> /dev/null)
	$(shell GO111MODULE=off go install github.com/boyter/scc@latest)
endif

.PHONY: clean_vendor
clean_vendor:
	rm -rf vendor go.sum

vendor:
	if [ ! -f go.mod ]; then go mod init; fi
	go mod tidy
	go mod vendor
	for cloudFunction in $(CLOUD_FUNCTIONS); do \
  		(cd cmd/functions/$$cloudFunction && go mod tidy) \
	done

.PHONY: revendor
revendor: clean_vendor vendor

.PHONY: clean_certs
clean_certs:
	rm -rf environments/local/certificates
	rm -rf environments/testing/certificates

.PHONY: certs
certs: clean_certs
	(mkdir -p environments/local/certificates && cd environments/local/certificates && mkcert -client -cert-file cert.pem -key-file key.pem api.prixfixe.local $(DEFAULT_CERT_TARGETS))
	(mkdir -p environments/testing/certificates && cd environments/testing/certificates && mkcert -client -cert-file cert.pem -key-file key.pem api.prixfixe.local $(DEFAULT_CERT_TARGETS))

## dependency injection

.PHONY: clean_wire
clean_wire:
	rm -f $(THIS)/internal/build/server/wire_gen.go

.PHONY: wire
wire: ensure_wire_installed vendor
	wire gen $(THIS)/internal/build/server

.PHONY: rewire
rewire: ensure_wire_installed clean_wire wire

## formatting

.PHONY: format
format: format_imports
	for file in `find $(PWD) -name '*.go'`; do $(GO_FORMAT) $$file; done

.PHONY: format_imports
format_imports:
	gci write --skip-generated --section standard --section default --section "prefix(github.com/prixfixeco)" --section "prefix(github.com/prixfixeco/backend)" .

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

.PHONY: docker_lint
docker_lint:
	@docker pull openpolicyagent/conftest:v0.28.3
	docker run --rm --volume $(PWD):$(PWD) --workdir=$(PWD) openpolicyagent/conftest:v0.21.0 test --policy docker_security.rego `find . -type f -name "*.Dockerfile"`

.PHONY: queries_lint
queries_lint:
	$(SQL_GENERATOR) compile

.PHONY: querier
querier: queries_lint
	$(SQL_GENERATOR) generate

.PHONY: golang_lint
golang_lint:
	@docker pull golangci/golangci-lint:v1.50
	docker run \
		--rm \
		--volume $(PWD):$(PWD) \
		--workdir=$(PWD) \
		golangci/golangci-lint:v1.50 golangci-lint run --config=.golangci.yml ./...

.PHONY: lint
lint: docker_lint queries_lint golang_lint # terraform_lint

.PHONY: clean_coverage
clean_coverage:
	@rm -f $(COVERAGE_OUT) profile.out;

.PHONY: coverage
coverage: clean_coverage $(ARTIFACTS_DIR)
	@go test -coverprofile=$(COVERAGE_OUT) -shuffle=on -covermode=atomic -race $(TESTABLE_PACKAGE_LIST) > /dev/null
	@go tool cover -func=$(ARTIFACTS_DIR)/coverage.out | grep 'total:' | xargs | awk '{ print "COVERAGE: " $$3 }'

.PHONY: quicktest # basically only running once instead of with -count 5 or whatever
quicktest: $(ARTIFACTS_DIR) vendor clear
	go build $(TOTAL_PACKAGE_LIST)
	go test -cover -shuffle=on -race -failfast $(TESTABLE_PACKAGE_LIST)

.PHONY: check_queries
check_queries:
	$(SQL_GENERATOR) compile

## Generated files

.PHONY: configs
configs:
	go run github.com/prixfixeco/backend/cmd/tools/gen_configs

.PHONY: queries
queries:
	go run github.com/prixfixeco/backend/cmd/tools/gen_queries

gen: configs queries

clean_ts:
	rm -rf $(ARTIFACTS_DIR)/typescript

typescript: clean_ts
	mkdir -p $(ARTIFACTS_DIR)/typescript
	go run github.com/prixfixeco/backend/cmd/tools/gen_typescript

## Integration tests

.PHONY: wipe_docker
wipe_docker:
	@docker stop $(shell docker ps -aq) && docker rm $(shell docker ps -aq)

.PHONY: docker_wipe
docker_wipe: wipe_docker

.PHONY: lintegration_tests # this is just a handy lil' helper I use sometimes
lintegration_tests: lint clear integration-tests

.PHONY: integration_tests
integration_tests: integration_tests_postgres

.PHONY: integration-tests
integration-tests: integration_tests_postgres

.PHONY: integration_tests_postgres
integration_tests_postgres:
	docker-compose \
	--file $(TEST_DOCKER_COMPOSE_FILES_DIR)/$(if $(filter y Y yes YES true TRUE plz sure yup YUP,$(OBSERVE)),integration-tests-with-observability.yaml,integration-tests.yaml) \
	up \
	--build \
	--quiet-pull \
	--force-recreate \
	--remove-orphans \
	--always-recreate-deps \
	$(if $(filter y Y yes YES true TRUE plz sure yup YUP,$(LET_HANG)),,--abort-on-container-exit) \
	--renew-anon-volumes

## Running

.PHONY: dev
dev: $(ARTIFACTS_DIR)
	docker-compose \
	--file $(ENVIRONMENTS_DIR)/local/compose_files/docker-compose.yaml up \
	--quiet-pull \
	--no-recreate \
	# --abort-on-container-exit \
	--always-recreate-deps

.PHONY: init_db
init_db: initialize_database

.PHONY: db_init
db_init: initialize_database

.PHONY: initialize_database
initialize_database:
	go run github.com/prixfixeco/backend/cmd/tools/db_initializer

## misc

.PHONY: tree
tree:
	tree -d -I vendor

.PHONY: line_count
line_count: ensure_scc_installed
	@scc --include-ext go --exclude-dir vendor

## maintenance

# https://cloud.google.com/sql/docs/postgres/connect-admin-proxy#connect-tcp
.PHONY: start_dev_cloud_sql_proxy
start_dev_cloud_sql_proxy:
	cloud_sql_proxy -dir=/cloudsql -instances='prixfixe-dev:us-central1:dev=tcp:5434'

.PHONY: proxy_dev_db
proxy_dev_db: start_dev_cloud_sql_proxy

dump_dev_db:
	rm -f cmd/tools/db_initializer/db_dumps/dump.sql
	for table in valid_preparations valid_measurement_units valid_instruments valid_ingredients valid_ingredient_preparations valid_ingredient_measurement_units valid_preparation_instruments recipes recipe_steps recipe_step_products recipe_step_instruments recipe_step_ingredients meals meal_recipes; do \
		pg_dump "user=api_db_user password=`gcloud secrets versions access --secret=api_user_database_password 1` host=127.0.0.1 port=5434 sslmode=disable dbname=prixfixe" --table="$$table" --data-only --column-inserts >> cmd/tools/db_initializer/dump.sql; \
	done
