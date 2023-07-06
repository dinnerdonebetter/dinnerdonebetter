PWD                           := $(shell pwd)
GOPATH                        := $(GOPATH)
ARTIFACTS_DIR                 := artifacts
COVERAGE_OUT                  := $(ARTIFACTS_DIR)/coverage.out
GO                            := docker run --interactive --tty --volume $(PWD):$(PWD) --workdir $(PWD) --user $(shell id -u):$(shell id -g) golang:1.18-stretch go
GO_FORMAT                     := gofmt -s -w
THIS                          := github.com/dinnerdonebetter/backend
TOTAL_PACKAGE_LIST            := `go list $(THIS)/...`
TESTABLE_PACKAGE_LIST         := `go list $(THIS)/... | grep -Ev '(integration)'`
ENVIRONMENTS_DIR              := environments
TEST_ENVIRONMENT_DIR          := $(ENVIRONMENTS_DIR)/testing
TEST_DOCKER_COMPOSE_FILES_DIR := $(TEST_ENVIRONMENT_DIR)/compose_files
SQL_GENERATOR                 := docker run --rm --volume `pwd`:/src --workdir /src kjconroy/sqlc:1.17.2
GENERATED_QUERIES_DIR         := internal/database/postgres/generated
LINTER_IMAGE                  := golangci/golangci-lint:v1.53.3
CONTAINER_LINTER_IMAGE        := openpolicyagent/conftest:v0.43.1
CLOUD_JOBS                    := meal_plan_finalizer meal_plan_grocery_list_initializer meal_plan_task_creator search_data_index_scheduler
CLOUD_FUNCTIONS               := data_changes outbound_emailer search_indexer
WIRE_TARGETS                  := server/http/build

## non-PHONY folders/files

regit:
	(rm -rf .git && cd ../ && rm -rf backend2 && git clone git@github.com:dinnerdonebetter/backend backend2 && cp -rf backend2/.git backend/.git && rm -rf backend2)

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

ensure_fieldalignment_installed:
ifndef $(shell command -v wire 2> /dev/null)
	$(shell GO111MODULE=off go get -u golang.org/x/tools/...)
endif

ensure_tagalign_installed:
ifndef $(shell command -v wire 2> /dev/null)
	$(shell go install github.com/4meepo/tagalign/cmd/tagalign@latest)
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
	for thing in $(CLOUD_FUNCTIONS); do \
  		(cd cmd/functions/$$thing && go mod tidy) \
	done
	for thing in $(CLOUD_JOBS); do \
  		(cd cmd/jobs/$$thing && go mod tidy) \
	done

.PHONY: revendor
revendor: clean_vendor vendor

## dependency injection

.PHONY: clean_wire
clean_wire:
	for tgt in $(WIRE_TARGETS); do \
		rm -f $(THIS)/internal/$$tgt/wire_gen.go; \
	done

.PHONY: wire
wire: ensure_wire_installed
	for tgt in $(WIRE_TARGETS); do \
		wire gen $(THIS)/internal/$$tgt; \
	done

.PHONY: rewire
rewire: clean_wire wire

## formatting

.PHONY: format
format: format_imports format_golang

.PHONY: format_golang
format_golang:
	@until fieldalignment -fix ./...; do true; done > /dev/null
	@until tagalign -fix -sort -order "json,toml" ./...; do true; done > /dev/null
	for file in `find $(PWD) -name '*.go'`; do $(GO_FORMAT) $$file; done

.PHONY: format_imports
format_imports:
	@# TODO: find some way to use $THIS here instead of hardcoding the path
	gci write --skip-generated --section standard --section "prefix(github.com/dinnerdonebetter/backend)" --section "prefix(github.com/dinnerdonebetter)" --section default --custom-order `find $(PWD) -type f -not -path '*/vendor/*' -name "*.go"`

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

.PHONY: lint_docker
lint_docker:
	@docker pull $(CONTAINER_LINTER_IMAGE)
	docker run --rm --volume $(PWD):$(PWD) --workdir=$(PWD) $(CONTAINER_LINTER_IMAGE) test --policy docker_security.rego `find . -type f -name "*.Dockerfile"`

.PHONY: queries_lint
queries_lint:
	$(SQL_GENERATOR) compile

.PHONY: querier
querier: queries_lint
	$(SQL_GENERATOR) generate

.PHONY: golang_lint
golang_lint:
	@docker pull $(LINTER_IMAGE)
	docker run \
		--rm \
		--volume $(PWD):$(PWD) \
		--workdir=$(PWD) \
		$(LINTER_IMAGE) golangci-lint run --config=.golangci.yml --timeout 15m ./...

.PHONY: lint
lint: lint_docker queries_lint golang_lint # terraform_lint

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
	go run github.com/dinnerdonebetter/backend/cmd/tools/gen_configs

.PHONY: queries
queries:
	go run github.com/dinnerdonebetter/backend/cmd/tools/gen_queries

gen: configs queries

clean_ts:
	rm -rf $(ARTIFACTS_DIR)/typescript

typescript: clean_ts
	mkdir -p $(ARTIFACTS_DIR)/typescript
	go run github.com/dinnerdonebetter/backend/cmd/tools/gen_clients/gen_typescript
	(cd ../frontend && make format)

clean_swift:
	rm -rf $(ARTIFACTS_DIR)/swift

swift: clean_swift
	mkdir -p $(ARTIFACTS_DIR)/swift
	go run github.com/dinnerdonebetter/backend/cmd/tools/codegen/gen_swift

## Integration tests

.PHONY: wipe_docker
wipe_docker:
	@docker stop $(shell docker ps -aq) && docker rm $(shell docker ps -aq)

.PHONY: docker_wipe
docker_wipe: wipe_docker

.PHONY: integration-tests
integration-tests: integration_tests

.PHONY: integration_tests
integration_tests: integration_tests_postgres

.PHONY: integration_tests_postgres
integration_tests_postgres:
	docker-compose \
	--file $(TEST_DOCKER_COMPOSE_FILES_DIR)/integration-tests.yaml \
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

.PHONY: start_infra
start_infra:
	docker-compose \
	--file $(ENVIRONMENTS_DIR)/local/compose_files/docker-compose.yaml up \
	--detach \
	--remove-orphans \
	postgres worker_queue

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
	cloud_sql_proxy -dir=/cloudsql -instances='dinner-done-better-dev:us-central1:dev=tcp:5434'

.PHONY: proxy_dev_db
proxy_dev_db: start_dev_cloud_sql_proxy

.PHONY: twirp  # I'm right here
twirp: clean_protobufs protobufs_typescript protobufs_golang

.PHONY: clean_protobufs
clean_protobufs:
	rm -rf internal/proto

.PHONY: protobufs_golang
protobufs_golang: $(ARTIFACTS_DIR)
	protoc --go_out=$(ARTIFACTS_DIR)/ \
		--twirp_out=$(ARTIFACTS_DIR)/ \
		--go_opt=paths=import \
		--proto_path protobufs \
		--experimental_allow_proto3_optional \
	 	protobufs/*.proto
	mv $(ARTIFACTS_DIR)/$(THIS)/internal/proto internal/
	rm -rf $(ARTIFACTS_DIR)/github.com

## Required for the typescript portions of this to work:
# npm install -g ts-protobufs twirp-ts @protobuf-ts/plugin@next

.PHONY: protobufs_typescript
protobufs_typescript: $(ARTIFACTS_DIR)
	protoc \
		-I ./protobufs \
		--plugin=protoc-gen-ts_proto=`which protoc-gen-ts_proto` \
		--plugin=protoc-gen-twirp_ts=`which protoc-gen-twirp_ts` \
		--ts_proto_opt=esModuleInterop=true \
		--ts_proto_opt=outputClientImpl=false \
		--ts_proto_out=$(ARTIFACTS_DIR) \
		--twirp_ts_opt="ts_proto" \
		--twirp_ts_out=$(ARTIFACTS_DIR) \
		--experimental_allow_proto3_optional \
		./protobufs/*.proto
