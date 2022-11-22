# we will put our integration testing in this path
INTEGRATION_TEST_PATH?=./integration_tests

# set of env variables that you need for testing
POSTGRES_URI=postgresql://postgres:test@localhost:5432/postgres

# this command will start a docker components that we set in docker-compose.yml
docker.start.components:
	docker-compose -f docker-compose-test.yaml up -d --remove-orphans postgres;

# shutting down docker components
docker.stop:
	docker-compose -f docker-compose-test.yaml down -v

# this command will trigger integration test
# INTEGRATION_TEST_SUITE_PATH is used for run specific test in Golang, if it's not specified
# it will run all tests under ./integration_tests directory
test.integration:
	$(MAKE) docker.start.components
	- go test -tags=integration $(INTEGRATION_TEST_PATH) -count=1 -run=$(INTEGRATION_TEST_SUITE_PATH) --integration --postgres-uri=$(POSTGRES_URI)
	$(MAKE) docker.stop

# this command will trigger integration test with verbose mode
test.integration.debug:
	$(MAKE) docker.start.components
	- go test -tags=integration $(INTEGRATION_TEST_PATH) -count=1 -v -run=$(INTEGRATION_TEST_SUITE_PATH) --integration --postgres-uri=$(POSTGRES_URI)
	$(MAKE) docker.stop
