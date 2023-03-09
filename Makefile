postgres:
		# docker run --name postgres14 -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -d postgres:14-alpine
		docker run --name postgres14 --network bank-network -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -d postgres:14-alpine

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

migrateup1:
		migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5432/simplebank?sslmode=disable" -verbose up 1

migratedown:
		migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5432/simplebank?sslmode=disable" -verbose down

migratedown1:
		migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5432/simplebank?sslmode=disable" -verbose down 1

sqlc:
		sqlc generate

test:
		go test -v -cover ./...

# shuffletest:
# 		go test -shuffle=on ./...

db_docs:
		dbdocs build docs/db.dbml

db_schema:
		dbml2sql --postgresql -o docs/schema.sql docs/db.dbml

server:
		go run main.go
	
mock:
		mockgen -package mockdb -destination db/mock/store.go github.com/blessedmadukoma/go-simple-bank/db/sqlc Store

proto:
		rm -f pb/*.go
		protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    proto/*.proto

.PHONY: postgres createdb dropdb psql createmigration migrateup migratedown migrateup1 migratedown1 sqlc test db_docs db_schema server mock proto