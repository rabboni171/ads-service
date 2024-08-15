run: run-ads-service run-auth-service run-metrics

add-docker-network:
	docker network create ads-network

run-ads-service:
	docker-compose -f ads-service/docker-compose.yml up -d

run-auth-service:
	docker-compose -f auth-service/docker-compose.yml up -d

run-metrics:
	docker-compose -f collecting-metrics/docker-compose.yml up -d
	