postgresinit:
	sudo docker run --name postgres -p 5435:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres

createdb:
	sudo docker exec -it postgres createdb --username=root --owner=root go-chat

postgres:
	sudo docker exec -it postgres psql

dropdb:
	sudo docker exec -it postgres dropdb go-chat

migrateup:
	migrate -path db/migrations -database "postgresql://root:password@localhost:5435/go-chat?sslmode=disable" -verbose up
	
migratedown:
	migrate -path db/migrations -database "postgresql://root:password@localhost:5435/go-chat?sslmode=disable" -verbose down

localstackinit:
	sudo docker run --rm -it -p 4566:4566 -p 4510-4559:4510-4559 localstack/localstack

.PHONY: postgresinit createdb postgres dropdb migrateup migratedown localstackinit