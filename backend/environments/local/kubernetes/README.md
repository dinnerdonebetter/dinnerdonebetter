# Kubernetes

Here lie the artifacts of trying to get the service running in kubernetes locally. I learned a lot, namely that I should probably just use Helm, but this will be handy to look back upon.

## Make targets

```
## EXPERIMENTAL KUBERNETES ZONE

.PHONY: k9s
k9s:
	k9s --refresh=1

.PHONY: build_localdev_api_container
build_localdev_api_container:
	docker build --tag dinner-done-better-api-server:latest --file environments/local/server.Dockerfile .

.PHONY: destroy_k8s
destroy_k8s:
	kubectl delete namespace localdev

.PHONY: deploy_to_k8s
deploy_to_k8s: build_localdev_api_container
	kubectl apply \
		--filename environments/local/kubernetes/0_namespace.yaml \
		--filename environments/local/kubernetes/1_postgres.yaml
	$(MAKE) wait_a_bit
	kubectl apply --filename environments/local/kubernetes/2_service.yaml

.PHONY: port_forward_db
port_forward_db:
	kubectl port-forward deployment/postgres 5434:5432 --namespace localdev

.PHONY: postgres_test
postgres_test:
	kubectl run --namespace localdev -it --rm psql-client --image=postgres:17 -- psql -W -h postgres -U dbuser -d dinner-done-better

.PHONY: wait_a_sec
wait_a_sec:
	sleep 1

.PHONY: wait_a_bit
wait_a_bit:
	sleep 15

.PHONY: wait_a_minute
wait_a_minute:
	sleep 60
```