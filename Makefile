project=loginsx
up:
	docker compose up
upd:
	docker compose up -d
upbuild:
	docker compose up -d --build
up-prod: 
	docker compose -f docker-compose.prod.yaml up --build -d
down:
	docker compose down
logs:
	docker logs ${project}
logs-follow:
	docker logs --follow ${project}
clear-logs:
	echo "" > $(docker inspect --format='{{.LogPath}}' ${project})