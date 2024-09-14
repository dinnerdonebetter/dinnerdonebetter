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
	echo "cd ../"
	echo "git clone git@github.com:dinnerdonebetter/dinnerdonebetter tempdir"
	@if [ -n "$(BRANCH)" ]; then \
	  echo "(cd tempdir && git checkout $(BRANCH))"; \
	fi
	echo "cp -rf tempdir/.git dinnerdonebetter/"
	echo "rm -rf tempdir"
