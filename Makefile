build-backend:
	docker compose -f docker-compose.yml build

start-backend:
	 docker-compose -f docker-compose.yml up tasks db

stop-all:
	docker-compose -f docker-compose.yml down -v

start-grafana:
	docker-compose -f docker-compose.yml up -d influxdb grafana

stop-grafana:
	docker-compose -f docker-compose.yml stop influxdb grafana

run-k6:
	docker-compose run k6 run \
		--env HOSTNAME=http://tasks:1323 \
		--vus 30 \
		--duration 30s \
		/scripts/script.js
#		--http-debug="full" \
