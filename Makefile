PWD      := $(shell pwd)
MYSELF   := $(shell id -u)
MY_GROUP := $(shell id -g)

.PHONY: setup
setup:
	(cd backend && $(MAKE) setup)
	(cd frontend && $(MAKE) setup)

.PHONY: format
format:
	(cd backend && $(MAKE) format)
	(cd frontend && $(MAKE) format)

.PHONY: lint
lint:
	(cd backend && $(MAKE) lint)
	(cd frontend && $(MAKE) lint)

.PHONY: test
test:
	(cd backend && $(MAKE) test)
	(cd frontend && $(MAKE) test)

.PHONY: openapi-clients
openapi-clients:
	(cd backend && $(MAKE) openapi-client)
	(cd frontend && $(MAKE) openapi-client)

.PHONY: openapi-lint
openapi-lint:
	npx @stoplight/spectral-cli@v6.13.1 lint openapi_spec.yamls

.PHONY: regit
regit:
	cd ../
	git clone git@github.com:dinnerdonebetter/dinnerdonebetter tempdir
	@if [ -n "$(BRANCH)" ]; then \
	  (cd tempdir && git checkout $(BRANCH)); \
	fi
	cp -rf tempdir/.git .
	rm -rf tempdir

.PHONY: terraformat
terraformat:
	(cd backend && $(MAKE) terraformat)
	(cd frontend && $(MAKE) terraformat)

#### NEW DEV K8S ENVIRONMENT ZONE

DEV_NAMESPACE := dev

.PHONY: helm_deps
helm_deps:
	helm repo add jetstack https://charts.jetstack.io --force-update
	helm repo add external-dns https://kubernetes-sigs.github.io/external-dns/ --force-update

.PHONY: deploy_dev
deploy_dev: helm_deps
	# kubectl delete deployments,cronjobs,configmaps,services,secrets --namespace $(DEV_NAMESPACE) --selector='managed_by!=terraform'
	skaffold run --filename=skaffold.yaml --build-concurrency 3 --profile $(DEV_NAMESPACE)
