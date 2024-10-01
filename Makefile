SERVICE_NAME=postgres-test
COMPOSE_FILE=docker-compose.yml

.PHONY: run-db test

run-db:
	@if [ -z "$$(docker compose -f $(COMPOSE_FILE) ps -q $(SERVICE_NAME))" ]; then \
  		docker compose -f $(COMPOSE_FILE) up -d $(SERVICE_NAME); \
 	else \
  		echo "Сервис '$(SERVICE_NAME)' уже запущен."; \
 	fi

test:
	migrate -source file://test_migrations -database postgres://test-user:test-password@${POSTGRES_TEST_HOST}:${POSTGRES_TEST_PORT}/test-db?sslmode=disable force 20240806143730
	migrate -source file://test_migrations -database postgres://test-user:test-password@${POSTGRES_TEST_HOST}:${POSTGRES_TEST_PORT}/test-db?sslmode=disable down -all
	migrate -source file://test_migrations -database postgres://test-user:test-password@${POSTGRES_TEST_HOST}:${POSTGRES_TEST_PORT}/test-db?sslmode=disable up
	cd tests && go test . -tags=unit
	cd tests && go test . -tags=integration
	cd tests && go test . -tags=e2e
	/home/nastya/allure-2.30.0/bin/allure generate ./tests/allure-results --clean -o ./tests/allure-report
	mkdir -p ./tests/allure-results/history && cp -r ./tests/allure-report/history/* ./tests/allure-results/history/ || true
#	/home/nastya/allure-2.30.0/bin/allure serve ./tests/allure-results

test_coverage:
	go test ./tests -coverprofile=coverage.out


