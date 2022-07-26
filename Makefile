PWD                           := $(shell pwd)
GOPATH                        := $(GOPATH)
ARTIFACTS_DIR                 := artifacts
COVERAGE_OUT                  := $(ARTIFACTS_DIR)/coverage.out
GO                            := docker run --interactive --tty --volume $(PWD):$(PWD) --workdir $(PWD) --user $(shell id -u):$(shell id -g) golang:1.18-stretch go
GO_FORMAT                     := gofmt -s -w
THIS                          := github.com/prixfixeco/api_server
TOTAL_PACKAGE_LIST            := `go list $(THIS)/...`
TESTABLE_PACKAGE_LIST         := `go list $(THIS)/... | grep -Ev '(cmd|tests|testutil|mock|fake)'`
ENVIRONMENTS_DIR              := environments
TEST_ENVIRONMENT_DIR          := $(ENVIRONMENTS_DIR)/testing
TEST_DOCKER_COMPOSE_FILES_DIR := $(TEST_ENVIRONMENT_DIR)/compose_files
LOCAL_ADDRESS                 := api.prixfixe.local
DEFAULT_CERT_TARGETS          := $(LOCAL_ADDRESS) prixfixe.local localhost 127.0.0.1 ::1

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

.PHONY: configs
configs:
	go run cmd/tools/config_gen/main.go

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

ensure_pnpm_installed:
ifndef $(shell command -v pnpm 2> /dev/null)
	$(shell npm install -g pnpm)
endif

.PHONY: ensure_hosts
ensure_hosts:
	if [ `cat /etc/hosts | grep api.prixfixe.local | wc -l` -ne 1 ]; then sudo -- sh -c "echo \"127.0.0.1       api.prixfixe.local\" >> /etc/hosts"; fi
	

.PHONY: clean_vendor
clean_vendor:
	rm -rf vendor go.sum

vendor:
	if [ ! -f go.mod ]; then go mod init; fi
	go mod tidy
	go mod vendor

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
format:
	for file in `find $(PWD) -name '*.go'`; do $(GO_FORMAT) $$file; done

.PHONY: terraformat
terraformat:
	@touch environments/dev/terraform/service-config.json
	@touch environments/dev/terraform/worker-config.json
	@touch environments/dev/terraform/opentelemetry-config.yaml
	@(cd environments/dev/terraform && terraform fmt)

.PHONY: check_terraform
check_terraform:
	@(cd environments/dev/terraform && terraform init -upgrade && terraform validate && terraform fmt && terraform fmt -check)

.PHONY: fmt
fmt: format terraformat

.PHONY: check_formatting
check_formatting: vendor
	docker build --tag check_formatting --file environments/testing/dockerfiles/formatting.Dockerfile .
	docker run --interactive --tty --rm check_formatting

## Testing things

.PHONY: pre_lint
pre_lint:
	@until fieldalignment -fix ./...; do true; done > /dev/null
	@echo ""

.PHONY: docker_lint
docker_lint:
	@docker pull openpolicyagent/conftest:v0.28.3
	docker run --rm --volume $(PWD):$(PWD) --workdir=$(PWD) openpolicyagent/conftest:v0.21.0 test --policy docker_security.rego `find . -type f -name "*.Dockerfile"`

.PHONY: lint
lint: docker_lint # check_terraform
	@docker pull golangci/golangci-lint:v1.46.2
	docker run \
		--rm \
		--volume $(PWD):$(PWD) \
		--workdir=$(PWD) \
		golangci/golangci-lint:v1.46.2 golangci-lint run --config=.golangci.yml ./...

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

.PHONY: frontend_tests
frontend_tests:
	docker-compose --file $(TEST_DOCKER_COMPOSE_FILES_DIR)/frontend_tests.yaml up \
	--build \
	--quiet-pull \
	--force-recreate \
	--remove-orphans \
	--renew-anon-volumes \
	--always-recreate-deps \
	--abort-on-container-exit

## Generated files

clean_ts:
	rm -rf $(ARTIFACTS_DIR)/typescript

typescript: clean_ts
	mkdir -p $(ARTIFACTS_DIR)/typescript
	go run cmd/tools/gen_ts/main.go

## Integration tests

.PHONY: wipe_docker
wipe_docker:
	@docker stop $(shell docker ps -aq) && docker rm $(shell docker ps -aq)

.PHONY: docker_wipe
docker_wipe:
	@docker stop $(shell docker ps -aq) && docker rm $(shell docker ps -aq)

.PHONY: ensure_postgres_is_up
ensure_postgres_is_up:
	@echo "waiting for postgres"
	@until docker exec --interactive --tty postgres psql 'postgres://dbuser:hunter2@localhost:5432/prixfixe?sslmode=disable' -c 'SELECT 1'; do true; done > /dev/null

.PHONY: ensure_elasticsearch_is_up
ensure_elasticsearch_is_up:
	@echo "waiting for elasticsearch"
	@until docker run --interactive --tty --network=host curlimages/curl:7.79.1 curl --head "http://localhost:9200/*" ; do true; done > /dev/null

.PHONY: wipe_local_postgres
wipe_local_postgres: ensure_postgres_is_up
	@echo "wiping postgres"
	@until docker exec --interactive --tty postgres psql 'postgres://dbuser:hunter2@localhost:5432/prixfixe?sslmode=disable' -c 'DROP SCHEMA public CASCADE; CREATE SCHEMA public;'; do true; done > /dev/null

.PHONY: wipe_local_elasticsearch
wipe_local_elasticsearch:
	@echo "wiping elasticsearch"
	@docker run --interactive --tty --network=host curlimages/curl:7.79.1 curl -X DELETE "http://localhost:9200/*"

.PHONY: lintegration_tests # this is just a handy lil' helper I use sometimes
lintegration_tests: lint clear integration-tests

.PHONY: integration_tests
integration_tests: integration_tests_postgres

.PHONY: integration-tests
integration-tests: integration_tests_postgres

.PHONY: integration_tests_postgres
integration_tests_postgres:
	docker-compose \
	--file $(TEST_DOCKER_COMPOSE_FILES_DIR)/integration-tests.yaml \
	up \
	--build \
	--force-recreate \
	--remove-orphans \
	--renew-anon-volumes \
	$(if $(filter y Y yes YES true TRUE plz sure yup YUP,$(LET_HANG)),, --abort-on-container-exit) \
	--always-recreate-deps

## Running

.PHONY: dev
dev: $(ARTIFACTS_DIR)
	docker-compose \
	--file $(ENVIRONMENTS_DIR)/local/docker-compose.yaml up \
	--quiet-pull \
	--no-recreate \
	--always-recreate-deps

.PHONY: init_db
init_db: initialize_database

.PHONY: db_init
db_init: initialize_database

.PHONY: initialize_database
initialize_database:
	go run github.com/prixfixeco/api_server/cmd/tools/db_initializer

## misc

.PHONY: tree
tree:
	tree -d -I vendor

.PHONY: line_count
line_count: ensure_scc_installed
	@scc --include-ext go --exclude-dir vendor

## maintenance

.PHONY: start_%_cloud_sql_proxy
start_%_cloud_sql_proxy:
	cloud_sql_proxy -dir=/cloudsql -instances=`gcloud sql instances describe $* | grep connectionName | cut -d ' ' -f 2`=tcp:5434
