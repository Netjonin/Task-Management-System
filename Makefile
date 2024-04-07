run/api:
	go run ./cmd/api

db/psql:
	psql ${TMS_DB_DSN}

db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

db/migrations/up:
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${TMS_DB_DSN} up