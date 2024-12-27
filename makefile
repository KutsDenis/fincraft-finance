TEST_DB_DSN := postgres://test_user:test_password@localhost:5433/test_db?sslmode=disable
.PHONY: test generate-mocks generate-proto

PROTO_DIR := ./api
PROTO_OUT := ./api

# Генерация протофайлов
# Проверил только на Windows, надеюсь на Linux тоже будет работать :)
generate-proto:
ifeq ($(OS),Windows_NT)
	powershell -Command "$$protoPath = Resolve-Path $(PROTO_DIR); Get-ChildItem -Recurse $(PROTO_DIR) -Filter *.proto | ForEach-Object { protoc --proto_path=$$protoPath --go_out=$(PROTO_OUT) --go-grpc_out=$(PROTO_OUT) $$_.FullName.Replace('\\', '/') }"
else
	find $(PROTO_DIR) -name "*.proto" -exec protoc --proto_path=$(PROTO_DIR) --go_out=$(PROTO_OUT) --go-grpc_out=$(PROTO_OUT) {} +
endif

# Генерация моков
generate-mocks:
	go generate ./...

# Запуск тестов с генерацией моков и протофайлов
test: generate-mocks generate-proto
	go test ./...

# Запуск тестов с генерацией моков и протофайлов используя gotestsum
test-sum: generate-mocks generate-proto
	gotestsum --format short-verbose ./...

# Запуск тестов c генерацией моков, протофайлов и сбором покрытия
test-coverage: generate-mocks generate-proto
	@rm -f coverage.out
	@go test -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out