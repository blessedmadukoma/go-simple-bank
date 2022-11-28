postgres:
		docker run --name postgres14 -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -d postgres:14-alpine

createdb:
		docker exec -it postgres14 createdb --username=postgres --owner=postgres simplebank

dropdb:
		docker exec -it postgres14 dropdb --username=postgres simplebank

psql: # log in to simplebank db in psql terminal
		docker exec -it postgres14 psql -U postgres -d simplebank

createmigration:
		migrate -help
		migrate create -ext sql -dir db/migration -seq [ADDNAME]

migrateup:
		migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5432/simplebank?sslmode=disable" -verbose up

migratedown:
		migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5432/simplebank?sslmode=disable" -verbose down

sqlc:
		sqlc generate

test:
		go test -v -cover ./...

server:
		go run main.go

.PHONY: postgres createdb dropdb psql createmigration migrateup migratedown sqlc test