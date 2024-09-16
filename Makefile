SERVICE_NAME=postgres-test
COMPOSE_FILE=docker-compose.yml

.PHONY: run-db test

run-db:
	@if [ -z "$$(docker compose -f $(COMPOSE_FILE) ps -q $(SERVICE_NAME))" ]; then \
  		docker compose -f $(COMPOSE_FILE) up -d $(SERVICE_NAME); \
 	else \
  		echo "Сервис '$(SERVICE_NAME)' уже запущен."; \
 	fi

test: run-db
	go test ./tests
	/home/nastya/allure-2.30.0/bin/allure serve ./tests/allure-results

test_coverage:
	go test ./tests -coverprofile=coverage.out


