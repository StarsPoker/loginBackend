project=loginsx
up:
	docker compose up -d --build
up-prod: 
	docker compose -f docker-compose.prod.yaml up --build -d
down:
	docker compose down
logs:
	docker compose logs ${project}
logs-follow:
	docker compose logs --follow ${project}
clear-logs:
	echo "" > $(docker inspect --format='{{.LogPath}}' ${project})