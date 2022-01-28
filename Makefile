GOCMD=go
GOTEST=$(GOCMD) test
LINTER=golangci-lint
# we will put our integration testing in this path
INTEGRATION_TEST_PATH?=./internal/repository

# set of env variables that you need for testing
ENV_LOCAL_TEST=\
  POSTGRES_PASSWORD=postgres \
  POSTGRES_USER=postgres

.PHONY: run
run:
	docker compose up --build

.PHONY: lint
lint:
	$(LINTER) run

.PHONY: test.short
test.short:
	$(GOTEST) -short ./... -count=1

.PHONY: test.full
test.full: test.short test.integration

.PHONY: test.integration
test.integration: 
	docker compose -f docker-compose.integration.yml up -d --remove-orphans;
	$(ENV_LOCAL_TEST) \
	$(GOTEST) -tags=integration $(INTEGRATION_TEST_PATH) -count=1
	docker compose -f docker-compose.integration.yml down;

.PHONY: docker.start
docker.start:
	docker compose -f docker-compose.integration.yml up -d --remove-orphans;

.PHONY: docker.stop
docker.stop:
	docker compose -f docker-compose.integration.yml down;

.PHONY: integration.test
integration.test:
	$(ENV_LOCAL_TEST) \
	$(GOTEST) -tags=integration $(INTEGRATION_TEST_PATH) -count=1

.PHONY: integration.test.debug
integration.test.debug:
	$(ENV_LOCAL_TEST) \
	$(GOTEST) -tags=integration $(INTEGRATION_TEST_PATH) -count=1 -v