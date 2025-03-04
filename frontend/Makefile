NODE_PACKAGE_MANAGER = yarn
CONTAINER_RUNNER     = docker
APPS                 = admin web landing

.PHONY: regit
regit:
	(rm -rf .git && cd ../ && rm -rf frontend2 && git clone git@github.com:dinnerdonebetter/frontend frontend2 && cp -rf frontend2/.git frontend/.git && rm -rf frontend2)

.PHONY: clean
clean:
	for thing in $(APPS); do \
		(rm -rf node_modules apps/$$thing/node_modules apps/$$thing/.turbo apps/$$thing/.next) \
	done
	rm -rf node_modules \
		packages/eslint-custom-config/node_modules \
		packages/models/node_modules packages/**/.turbo

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
	@(cd deploy/environments/dev/webapp/terraform && terraform fmt)
	@(cd deploy/environments/dev/landing/terraform && terraform fmt)
	@(cd deploy/environments/dev/admin/terraform && terraform fmt)

.PHONY: lint_terraform
lint_terraform: terraformat
	@(cd deploy/environments/dev/admin/terraform && terraform init -upgrade && terraform validate && terraform fmt && terraform fmt -check)
	@(cd deploy/environments/dev/landing/terraform && terraform init -upgrade && terraform validate && terraform fmt && terraform fmt -check)
	@(cd deploy/environments/dev/webapp/terraform && terraform init -upgrade && terraform validate && terraform fmt && terraform fmt -check)

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

#### K8S ENVIRONMENT ZONE

DEV_NAMESPACE := dev

.PHONY: deploy_dev
deploy_dev:
	kubectl delete deployments --all --ignore-not-found=true --namespace $(DEV_NAMESPACE)
	skaffold run --filename=skaffold.yaml --build-concurrency 0 --profile $(DEV_NAMESPACE)
