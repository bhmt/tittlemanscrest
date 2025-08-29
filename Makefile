.DEFAULT_GOAL := unit

.PHONY: unit
unit:
	go test -v -cover -coverprofile=coverage_unit.out ./...

.PHONY: integ
integ: cdown cup
	go test -tags=integration -v -cover -coverprofile=coverage.out ./...

.PHONY: cup
cup:
	docker compose -f postgres.yaml up -d --remove-orphans

.PHONY: cdown
cdown:
	docker compose -f postgres.yaml down -v --remove-orphans

.PHONY: bench
bench:
	go test -benchmem -run=^$ -bench ^Benchmark github.com/bhmt/tittlemanscrest/cache
