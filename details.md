## Logs of what I did for each video

1. Lecture 1:
   1. Created database schema using dbdiagram.io
   
2. Lecture 2:
   1. Setup docker, and the postgres container in docker
     - Steps to run the postgres container: <br>
        1. Pull the image: <br> `docker pull <image>:<tag>` e.g. docker pull postgres:14-alpine
        2. List all images you've pulled: <br/> `docker images`
        3. Run a container from the image: <br/> `docker run --name <container_name> -e <env_variable> -p <host_ports> -d <image:tag>` e.g. docker run --name postgres14 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres:14-alpine
        4. List all running containers: <br/> `docker ps`
        5. Execute the container (run psql): <br/> `docker exec -it <container_name_or_id> <command> [args]` e.g. docker exec -it postgres14 psql -U postgres
   2. Setup the database by running the sql file.
   
3. Lecture 3:
   1. Installed golang-migrate for managing database migration
   2. Created a  `Makefile` to automate creation, dropping and running of db. To run any command in the `Makefile`: <br/> `make <command>`
   3. Create a migration file: <br> `migrate create -ext <extension_name> -dir <directory_location> -seq <file_name>` e.g. migrate create -ext sql -dir db/migration -seq init_schema
   4. Run `migrate` command to execute the migration: <br/> `migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5432/simplebank?sslmode=disable" -verbose up`
   
4. Lecture 4: Generate CRUD Golang code from SQL
   1. Installed `sqlc` using brew and set up the `sqlc.yaml` config.
   2. Added sqlc to Makefile for autogeneration of Go SQL code
   3. generated Accounts, Entries and Transfers struct and CRUD operation methods using sqlc

5. Lecture 5: Unit testing Go database CRUD
   1. installed [pq](github.com/lib/pq) database driver to set up the TestMain function which tests the database connection 
   2. installed [testify](github.com/stretchr/testify) to perform unit test operations for Creating, Getting, Updating and Deleting accounts, entries and transfers
   3. created util package for random generation of money, owner name, integers and currencies
   4. added Golang test command to the Makefile

6. Lecture 6: Golang DB Transaction
   1. Why we need db transaction:
      1. to provide a reliable and consistent unit of work, even in case of system failure
      2. to provide isolation between programs that access the database concurrently
   2. created method and store struct (involves composition) for Transaction which performs the following operations: create a transfer record, add account entries, update accounts' balance
   3. wrote test for transfer transaction
   
7. Lecture 7: How to Debug a Deadlock
   1. added transfer feature
   2. added code in sql to prevent deadlock
   3. update balance of both accounts
   
8. Lecture 8: How to avoid DB Deadlock
   1. added test for Deadlock
   2. implemented addMoney function for easy code readability

9. Lecture 9: Transaction Isolation Level
   1.  MySQL uses 4 isolation levels - read uncommited, read commited, repeatable read, serializable
   2.  Postgres uses 3 isolation levels - read commited (same with read uncommited - which works like read commited, not read uncommited in MySQL), repeatable read, serializable

10. Lecture 10: Github Actions in Go and Postgres
    1.  Set up package directory for Github actions (CI/CD): `.github/workflows/`
    2.  added postgres service to the `ci.yml` workflow
    3.  added golang-migrate installation, `make migrateup` (for migration) and `make test` (for test) in the steps of the workflow.
    4.  added cURL command to install golang-migrate, move the installed golang-migrate into the user's bin folder (using the | or pipe) and added a `which migrate` command to see if migrate was installed successfully

11. Lecture 11: Implement HTTP API in Go using Gin
    1.  installed gin using the `go get` command
    2.  created a new folder `api` and a `server.go` file to handle API and server requests
    3.  created an errorResponse method to handle error response
    4.  implemented a createAccountRequest struct to handle input values of a new user (balance by default, will be 0)
    5.  created a `main.go` to handle running of the server, which the method was created in the `server.go` file.
    6.  Fixed the major bug when testing the transfer of money from one account to another i.e. another instance of the `testDB` was being instantiated in the main_test.go (i.e. testDB, err := sql.Open), instead of using the global `testDB` variable (testDB, err = sql.Open)
    7.  updated Makefile by adding `server` command to run the api
    8.  created endpoints for testing the create account and get account, list accounts (using pagination) features.
    9.  updated `sqlc.yaml` to emit empty slices, which returns empty slice ([]) instead of nil or null when the record is empty