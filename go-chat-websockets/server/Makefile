server:
	go build -o bin/server ./cmd/
	./bin/server

postgresinit:
	sudo docker run --name postgres-db -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=nopassword -d postgres:alpine
	sleep 5

.PHONY: createdb postgresinit dropdb psql migrate-up migrate-down server

createdb:
	sudo docker exec -it postgres-db createdb --owner=root go-chat

dropdb:
	sudo docker exec -it postgres-db dropdb go-chat

psql:
	sudo docker exec -it postgres-db psql

migrate-up:
	migrate -path db/migrations -database "postgresql://root:nopassword@localhost:5432/go-chat?sslmode=disable" -verbose up

migrate-down:
	migrate -path db/migrations -database "postgresql://root:nopassword@localhost:5432/go-chat?sslmode=disable" -verbose down