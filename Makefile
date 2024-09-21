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

.PHONY: openapi
openapi:
	(cd backend && make openapi-client)
	npx openapi-typescript-codegen@0.29.0 --input openapi_spec.yaml --output artifacts/generated
	$(MAKE format)

.PHONY: openapi-lint
openapi-lint:
	npx @stoplight/spectral lint openapi_spec.yaml

.PHONY: regit
regit:
	cd ../
	git clone git@github.com:dinnerdonebetter/dinnerdonebetter tempdir
	@if [ -n "$(BRANCH)" ]; then \
	  (cd tempdir && git checkout $(BRANCH)); \
	fi
	cp -rf tempdir/.git .
	rm -rf tempdir
	cd dinnerdonebetter
