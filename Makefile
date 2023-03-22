rebuild-backend:
	cd task-service && go mod tidy && docker compose -f docker-compose.yml build

run-backend:
	cd task-service && docker-compose -f docker-compose.yml up

stop-backend:
	cd task-service && docker-compose -f docker-compose.yml down -v