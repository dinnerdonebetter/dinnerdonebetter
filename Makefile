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

LOCAL_DEV_NAMESPACE := localdev

.PHONY: k9s
k9s:
	k9s --refresh 1 --namespace $(LOCAL_DEV_NAMESPACE)

.PHONY: skbuild
skbuild:
	skaffold build --build-concurrency 0

.PHONY: skrender
skrender:
	rm -f environments/local/generated/kubernetes.yaml
	mkdir -p deploy/environments/local/generated/
	$(MAKE) deploy/environments/local/generated/kubernetes.yaml

.PHONY: helm_deps
helm_deps:
	helm repo add open-telemetry https://open-telemetry.github.io/opentelemetry-helm-charts
	helm repo add grafana https://grafana.github.io/helm-charts
	helm repo add prometheus https://prometheus-community.github.io/helm-charts
	helm repo update

.PHONY: dev
dev: helm_deps nuke_k8s skrender
	skaffold dev --build-concurrency 0 --profile $(LOCAL_DEV_NAMESPACE) --port-forward

deploy/environments/local/generated/kubernetes.yaml:
	skaffold render --profile $(LOCAL_DEV_NAMESPACE) --output deploy/environments/local/generated/kubernetes.yaml

.PHONY: nuke_k8s
nuke_k8s:
	kubectl delete namespace $(LOCAL_DEV_NAMESPACE) || true

.PHONY: proxy_k8s_db
proxy_k8s_db:
	kubectl port-forward --namespace localdev svc/postgres-postgresql 5434:5432

.PHONY: proxy_k8s_api_server
proxy_k8s_api_server:
	kubectl port-forward --namespace localdev svc/dinner-done-better-api-svc 8888:8000

.PHONY: proxy_k8s_webapp_server
proxy_k8s_webapp_server:
	kubectl port-forward --namespace localdev svc/dinner-done-better-webapp-srv 9999:9000

.PHONY: proxy_k8s_admin-app_server
proxy_k8s_admin_server:
	kubectl port-forward --namespace localdev svc/dinner-done-better-admin-app-srv 7777:7000

.PHONY: proxy_k8s_landing_server
proxy_k8s_landing_server:
	kubectl port-forward --namespace localdev svc/dinner-done-better-landing-srv 60006:60000
