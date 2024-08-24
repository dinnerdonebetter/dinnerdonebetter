setup:
	(cd backend && make setup)
	(cd frontend && make setup)

format:
	(cd backend && make format)
	(cd frontend && make format)

.PHONY: regit
regit:
	(cd ../ && git clone git@github.com:dinnerdonebetter/dinnerdonebetter tempdir && (cd tempdir && git checkout go-1.23) && cp -rf tempdir/.git dinnerdonebetter/ && rm -rf tempdir)
