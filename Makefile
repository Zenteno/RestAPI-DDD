# all our targets are phony (no files to check).
.PHONY: shell help build rebuild service login test clean prune


# Regular Makefile part for buildpypi itself
help:
	@echo ''
	@echo 'Usage: make [TARGET] [EXTRA_ARGUMENTS]'
	@echo 'Targets:'
	@echo '  rebuild  	docker-compose build --no-cache'
	@echo '  test     	test docker --container-- '
	@echo '  up   		run as service docker-compose up -d '
	@echo ''
rebuild:
	# force a rebuild by passing --no-cache
	docker-compose build --no-cache

up:
	# run as a (background) service
	docker-compose up -d 

test:
	# here it is useful to add your own customised tests
	docker build -f DockerTest.Dockerfile -t golang_test .
	rm -rf coverage
	mkdir coverage
	docker run --rm --network $(shell basename $(CURDIR))_app-network -v $(CURDIR)/coverage:/app/coverage -it golang_test /bin/bash -c "/app/api-ddd.test -test.coverprofile=system.out && go tool cover -html=system.out -o=/app/coverage/index.html"
	mv coverage/index.html coverage.html
	rm -rf coverage
	docker rmi golang_test		