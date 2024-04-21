postgres:
		# docker run --name postgres14 -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -d postgres:14-alpine
		docker run --name postgres14 --network bank-network -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -d postgres:14-alpine

sbrunonce: # run simplebank docker image once
		docker run --name simple-bank-api --network bank-network -p 8080:8080 -e DB_URL=postgres://postgres:postgres@postgres14:5432/simplebank?sslmode=disable -d simple-bank-api

pgstart:
		docker start postgres14

pgstop:
		docker stop postgres14

sbstart:
		docker start simple-bank-api

sbstop:
		docker stop simple-bank-api

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

new_migration:
		migrate create -ext sql -dir db/migration -seq $(name)

sqlc:
		sqlc generate

test:
		go test -v -cover -short ./...

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
		mockgen -package mockwk -destination worker/mock/distributor.go github.com/blessedmadukoma/go-simple-bank/worker TaskDistributor

proto:
		rm -f pb/*.go
		rm -f docs/swagger/*.swagger.json
		protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
				--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
				--openapiv2_out=docs/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
    proto/*.proto
		statik -src=./docs/swagger -dest=./docs

evans:
		evans --host localhost --port 9090 -r repl

redis:
	docker run --name redis -p 6379:6379 -d redis:7.2-alpine

redisping:
	docker exec -it redis redis-cli ping

.PHONY: postgres createdb dropdb psql createmigration migrateup migratedown migrateup1 migratedown1 sqlc test db_docs db_schema server mock proto evans redis redisping new_migration