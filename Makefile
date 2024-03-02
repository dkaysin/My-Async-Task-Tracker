up_pg_docker:
	docker compose --profile pg_docker up --build

up_pg_local:
	docker compose --profile pg_local up --build

up_kafka:
	docker compose --profile kafka up --build

up_services:
	docker compose --profile services up --build

avro:
	cd schema_registry && make avro && cd ..
