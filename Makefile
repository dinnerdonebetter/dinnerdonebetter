.PHONY: setup
setup:
	(cd backend && make setup)
	(cd frontend && make setup)

.PHONY: format
format:
	(cd backend && make format)
	(cd frontend && make format)

.PHONY: regit
regit:
	(cd ../ && git clone git@github.com:dinnerdonebetter/dinnerdonebetter tempdir && cp -rf tempdir/.git dinnerdonebetter/ && rm -rf tempdir)
