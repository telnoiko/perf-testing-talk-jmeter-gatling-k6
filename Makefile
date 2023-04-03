build-backend:
	docker compose -f docker-compose.yml build

start-backend:
	 docker-compose -f docker-compose.yml up tasks

stop-all:
	docker-compose -f docker-compose.yml down -v

run-gatling:
	cd gatling && ./gradlew gatlingRun-tasks.TaskSimulation

start-grafana:
	docker-compose -f docker-compose.yml up -d influxdb grafana

run-k6:
	docker-compose run k6 run \
		--env HOSTNAME=http://tasks:1323 \
		--vus 30 \
		--duration 30s \
		/scripts/script.js
#		--http-debug="full" \

convert-postman-k6:
	postman-to-k6 ./postman/collection.json -e ./postman/environment.json -s -o ./k6/generated/script.js

run-k6-converted:
	docker-compose run k6 run \
		--vus 2 \
		--duration 30s \
		/scripts/generated/script.js
#		--http-debug="full" \
