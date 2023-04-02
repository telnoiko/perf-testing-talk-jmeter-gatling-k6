build-backend:
	docker compose -f docker-compose.yml build

start-backend:
	 docker-compose -f docker-compose.yml up tasks

stop-all:
	docker-compose -f docker-compose.yml down -v

run-gatling:
	cd gatling && ./gradlew gatlingRun-tasks.TaskSimulation

run-k6:
	docker-compose run k6 run \
		--env HOSTNAME=http://tasks:1323 \
		--vus 30 \
		--duration 30s \
		/scripts/script.js
#		--http-debug="full" \
