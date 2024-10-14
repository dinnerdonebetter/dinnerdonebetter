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
