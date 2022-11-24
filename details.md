## Logs of what I did for each video

1. Lesson 1:
   1. Created database schema using dbdiagram.io
   
2. Lesson 2:
   1. Setup docker, and the postgres container in docker
     - Steps to run the postgres container: <br>
        1. Pull the image: <br> `docker pull <image>:<tag>` e.g. docker pull postgres:14-alpine
        2. List all images you've pulled: <br/> `docker images`
        3. Run a container from the image: <br/> `docker run --name <container_name> -e <env_variable> -p <host_ports> -d <image:tag>` e.g. docker run --name postgres14 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres:14-alpine
        4. List all running containers: <br/> `docker ps`
        5. Execute the container (run psql): <br/> `docker exec -it <container_name_or_id> <command> [args]` e.g. docker exec -it postgres14 psql -U postgres
   2. Setup the database by running the sql file.
   
3. Lesson 3:
   1. Installed golang-migrate for managing database migration
   2. Created a  `Makefile` to automate creation, dropping and running of db. To run any command in the `Makefile`: <br/> `make <command>`
   3. Create a migration file: <br> `migrate create -ext <extension_name> -dir <directory_location> -seq <file_name>` e.g. migrate create -ext sql -dir db/migration -seq init_schema
   4. Run `migrate` command to execute the migration: <br/> `migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5432/simplebank?sslmode=disable" -verbose up`
   
4. Lesson 4: Generate CRUD Golang code from SQL
   1. Installed `sqlc` using brew and set up the `sqlc.yaml` config.
   2. Added sqlc to Makefile for autogeneration of Go SQL code
   3. generated Accounts, Entries and Transfers struct and CRUD operation methods using sqlc

5. Lesson 5: Unit testing Go database CRUD
   1. installed [pq](github.com/lib/pq) database driver to set up the TestMain function which tests the database connection 
   2. installed [testify](github.com/stretchr/testify) to perform unit test operations for Creating, Getting, Updating and Deleting accounts, entries and transfers
   3. created util package for random generation of money, owner name, integers and currencies
   4. added Golang test command to the Makefile

6. Lesson 6: Golang DB Transaction
   1. Why we need db transaction:
      1. to provide a reliable and consistent unit of work, even in case of system failure
      2. to provide isolation between programs that access the database concurrently
   2. created method and store struct (involves composition) for Transaction which performs the following operations: create a transfer record, add account entries, update accounts' balance
   3. wrote test for transfer transaction
   
7. Lesson 6: How to Debug a Deadlock