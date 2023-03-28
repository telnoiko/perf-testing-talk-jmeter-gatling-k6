build-backend:
	#cd task-service && go mod tidy
	docker compose -f docker-compose.yml build

start-backend:
	 docker-compose -f docker-compose.yml up

stop-backend:
	docker-compose -f docker-compose.yml down -v

run-k6:
	cd k6 && k6 run --env HOSTNAME=http://localhost:1323 --vus 10 --duration 30s script.js