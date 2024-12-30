PWD      := $(shell pwd)
MYSELF   := $(shell id -u)
MY_GROUP := $(shell id -g)

.PHONY: setup
setup:
	(cd backend && make setup)
	(cd frontend && make setup)

.PHONY: format
format:
	(cd backend && make format)
	(cd frontend && make format)

.PHONY: lint
lint:
	(cd backend && make lint)
	(cd frontend && make lint)

.PHONY: test
test:
	(cd backend && make test)
	(cd frontend && make test)

.PHONY: openapi-clients
openapi-clients:
	(cd backend && make openapi-client)
	(cd frontend && make openapi-client)

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


### EXPERIMENTAL KUBERNETES ZONE

DEV_NAMESPACE := dev
DEV_GENERATED_K8S := deploy/environments/dev/generated/kubernetes.yaml

.PHONY: k9s
k9s:
	k9s --refresh 1 --namespace $(DEV_NAMESPACE)

.PHONY: build
build:
	skaffold build --build-concurrency 0 --profile $(DEV_NAMESPACE)

.PHONY: skrender
skrender: clean_k8s
	mkdir -p deploy/environments/local/generated/
	$(MAKE) generate_kubernetes

.PHONY: dev
dev: helm_deps nuke_k8s skrender
	skaffold dev --build-concurrency 0 --profile $(DEV_NAMESPACE) --port-forward --tail=false

.PHONY: clean_k8s
clean_k8s:
	rm -f $(DEV_GENERATED_K8S)

.PHONY: generate_kubernetes
generate_kubernetes:
	skaffold render --profile $(DEV_NAMESPACE) --output $(DEV_GENERATED_K8S)

.PHONY: nuke_k8s
nuke_k8s:
	kubectl delete namespace $(DEV_NAMESPACE) || true
