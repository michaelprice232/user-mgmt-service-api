start-local-db:
	docker-compose up -d

run-webserver: start-local-db
	LOG_LEVEL=debug RUNNING_LOCALLY=true go run cmd/main.go

stop-database-and-delete-volume:
	docker-compose stop
	docker-compose rm --force
	docker volume rm user-mgmt-service-api_db-data --force

test-endpoints:
	./scripts/client.sh

prune-docker:
	docker system prune --all --volumes --force