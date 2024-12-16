SERVICE_NAME=postgres-test
COMPOSE_FILE=docker-compose.yml

.PHONY: run-db test auth

run-db:
	@if [ -z "$$(docker compose -f $(COMPOSE_FILE) ps -q $(SERVICE_NAME))" ]; then \
  		docker compose -f $(COMPOSE_FILE) up -d $(SERVICE_NAME); \
 	else \
  		echo "Сервис '$(SERVICE_NAME)' уже запущен."; \
 	fi

test:
	export POSTGRES_TEST_HOST=postgres-test POSTGRES_TEST_PORT=5431 
	migrate -source file://test_migrations -database postgres://test-user:test-password@${POSTGRES_TEST_HOST}:${POSTGRES_TEST_PORT}/test-db?sslmode=disable force 20240806143730
	migrate -source file://test_migrations -database postgres://test-user:test-password@${POSTGRES_TEST_HOST}:${POSTGRES_TEST_PORT}/test-db?sslmode=disable down -all
	migrate -source file://test_migrations -database postgres://test-user:test-password@${POSTGRES_TEST_HOST}:${POSTGRES_TEST_PORT}/test-db?sslmode=disable up
	cd tests && go test . -tags=unit
	cd tests && go test . -tags=integration
	cd tests && go test . -tags=e2e
	allure generate ./tests/allure-results --clean -o ./tests/allure-report
	mkdir -p ./tests/allure-results/history && cp -r ./tests/allure-report/history/* ./tests/allure-results/history/ || true
	@if [ $(ALLURE_SERVE) -eq "1" ]; then allure serve ./tests/allure-results; fi

auth:
	cd tests && go test . -tags=auth

test_coverage:
	go test ./tests -coverprofile=coverage.out


