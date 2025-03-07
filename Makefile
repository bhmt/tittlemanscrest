.PHONY: unit integ cdown cup
.DEFAULT_GOAL := unit

unit:
	go test -v -cover -coverprofile=coverage_unit.out ./...

integ:
	docker compose -f postgres.yaml down -v --remove-orphans || echo "compose down not run"
	docker compose -f postgres.yaml up -d --remove-orphans
	go test -tags=integration -v -cover -coverprofile=coverage.out ./...

cup:
	docker compose -f postgres.yaml up -d --remove-orphans

cdown:
	docker compose -f postgres.yaml down -v --remove-orphans
