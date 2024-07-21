ifneq (,$(wildcard ./.env))
    include .env
    export
endif

migrate-up:
	goose -dir db/migrations postgres "host=localhost port=${POSTGRES_PORT} user=${POSTGRES_USER} \
		password=${POSTGRES_PASSWORD} dbname=${POSTGRES_DB} sslmode=disable" up

migrate-down:
	goose -dir db/migrations postgres "host=localhost port=${POSTGRES_PORT} user=${POSTGRES_USER} \
		password=${POSTGRES_PASSWORD} dbname=${POSTGRES_DB} sslmode=disable" down

