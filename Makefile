binaries:
	@cd ./provider/aws && make
	@cd ./provider/gcp && make
	@cd ./provider/azure && make

check: lint test

fmt:
	@cd ./core && make fmt
	@cd ./provider/aws && make fmt
	@cd ./provider/gcp && make fmt
	@cd ./provider/azure && make fmt

lint:
	@cd ./core && make lint
	@cd ./provider/aws && make lint
	@cd ./provider/gcp && make lint
	@cd ./provider/azure && make lint

test:
	@cd ./core && make test
	@cd ./provider/aws && make test
	@cd ./provider/gcp && make test
	@cd ./provider/azure && make test

test-coverage:
	@cd ./core && make test
	@cd ./provider/aws && make test-coverage
	@cd ./provider/gcp && make test-coverage
	@cd ./provider/azure && make test-coverage

license-check:
	@cd ./provider/aws && make license-check
	@cd ./provider/gcp && make license-check
	@cd ./provider/azure && make license-check

generate-sources:
	@cd ./core && make generate-sources
	@cd ./provider/aws && make generate-sources
	@cd ./provider/gcp && make generate-sources
	@cd ./provider/azure && make generate-sources