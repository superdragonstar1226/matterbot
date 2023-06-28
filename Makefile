

DOCKER_COMPOSE=sudo docker-compose -f docker-compose.yaml

.PHONY: run docker-run docker-stop create-db drop-db gofmt

#  running mattermost-bot server.
run: 
	go run ./mattermost/cmd/mattermost/main.go

docker-run: 
	${DOCKER_COMPOSE} up -d

docker-stop:
	${DOCKER_COMPOSE} down

# gofmt entire project.
gofmt:
	gofmt -s -w .

create-db:
	@echo " Note: it may produce 'already exists' errors doesn't matter"
	PGPASSWORD=123456 psql -h localhost -U postgres \
		-f pkg/scripts/create_dbs.sql

drop-db:
	@echo " Note: it may produce 'already exists' errors doesn't matter"
	PGPASSWORD=123456 psql -h localhost -U postgres \
		-f pkg/scripts/drop_dbs.sql

migrateup:
	migrate -path pkg/scripts/migrations -database "postgresql://mattermost:secret@localhost:5432/mattermost?sslmode=disable" up
migratedown:
	migrate -path pkg/scripts/migrations -database "postgresql://mattermost:secret@localhost:5432/mattermost?sslmode=disable" down




