NODE_PACKAGE_MANAGER = yarn
CONTAINER_RUNNER     = docker
APPS                 = admin web landing

regit:
	(rm -rf .git && cd ../ && rm -rf frontend2 && git clone git@github.com:dinnerdonebetter/frontend frontend2 && cp -rf frontend2/.git frontend/.git && rm -rf frontend2)

clean:
	rm -rf node_modules \
	apps/web/node_modules apps/web/.turbo apps/web/.next \
	packages/eslint-custom-config/node_modules \
	packages/models/node_modules packages/models/.turbo

.PHONY: setup
setup:
	$(NODE_PACKAGE_MANAGER) install
	for thing in $(APPS); do \
  		(cd apps/$$thing && $(NODE_PACKAGE_MANAGER) install) \
	done

.PHONY: lint
lint:
	$(NODE_PACKAGE_MANAGER) lint

.PHONY: terraformat
terraformat:
	@(cd environments/dev/webapp/terraform && terraform fmt)
	@(cd environments/dev/landing/terraform && terraform fmt)
	@(cd environments/dev/admin/terraform && terraform fmt)

.PHONY: lint_terraform
lint_terraform: terraformat
	@(cd environments/dev/admin/terraform && terraform init -upgrade && terraform validate && terraform fmt && terraform fmt -check)
	@(cd environments/dev/landing/terraform && terraform init -upgrade && terraform validate && terraform fmt && terraform fmt -check)
	@(cd environments/dev/webapp/terraform && terraform init -upgrade && terraform validate && terraform fmt && terraform fmt -check)

.PHONY: dev
dev:
	$(NODE_PACKAGE_MANAGER) dev

.PHONY: stop
stop:
	kill -KILL `lsof -i tcp:9000 -i tcp:7000 | tail -n +2 | awk '{print $$2}'`

.PHONY: build
build:
	$(NODE_PACKAGE_MANAGER) build

.PHONY: test
test:
	$(NODE_PACKAGE_MANAGER) test

.PHONY: format
format:
	$(NODE_PACKAGE_MANAGER) format

.PHONY: format-check
format-check:
	$(NODE_PACKAGE_MANAGER) format-check

.PHONY: build_containers
build_containers: build_admin_container build_landing_container build_web_container
	@echo ""

.PHONY: build_admin_container
build_admin_container:
	$(CONTAINER_RUNNER) build --tag ddb_admin:latest --file apps/admin/Dockerfile .

.PHONY: build_landing_container
build_landing_container:
	$(CONTAINER_RUNNER) build --tag ddb_landing:latest --file apps/landing/Dockerfile .

.PHONY: build_web_container
build_web_container:
	$(CONTAINER_RUNNER) build --tag ddb_web:latest --file apps/web/Dockerfile .

.PHONY: openapi-client
openapi-client:
	rm -rf packages/api-client/generated
	npx openapi-typescript-codegen --input ../openapi_spec.yaml --output packages/api-client/generated --useUnionTypes --client fetch
	$(MAKE) format
