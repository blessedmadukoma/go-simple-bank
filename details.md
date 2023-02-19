## Logs of what I did

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
    7.  updated Makefile by adding `server` command to run the api.
    8.  created endpoints for testing the create account and get account, list accounts (using pagination) and update accounts features.
    9.  updated `sqlc.yaml` to emit empty slices, which returns empty slice ([]) instead of nil or null when the record is empty.
    
12. Lecture 12: Load config from file & env vars with Viper
    1.  installed [`viper`](github.com/spf13/viper) using the go get commmand.
    2.  created `app.env` to house the environmental variables
    3.  implemented LoadConfig() method in  `config.go` to load the env variables from the app.env file using viper.
    4.  updated main.go and main_test.go files to read from the LoadConfig() method.

13. Lecture 13: Mock database for testing HTTP API in Go
    1.  Why Mock DB?
        1.  independent tests - isolate test data to avoid conflicts
        2.  faster tests - reduce a lot of time talking to the db
        3.  100% coverage - easily setup edge cases such as unexpected cases
    2. installed [mockgen](github.com/golang/mock/mockgen@v1.6.0), a package in gomock
    3. create a store interface and rename Store struct to SQLStore struct which executes all SQL queries and transactions
    4. update Store interface by adding the TransferTx method
    5. update sqlc config (sqlc.yaml) -> emit_interface: true, then rerun `sqlc generate` or `make sqlc` which creates a Querier file
 14. update Store interface by adding the newly generated Querier interface, which holds all methods implementing the *Queries struct
    1. update server.go by removing the pointers from `store db.Store` on both struct and method parameters since Store is no longer a struct but an interface
    2. create a new folder `mock` in the db folder
    3. run a mockgen command on the module name which holds the sqlc queries and the interface, and add the destination (the mock folder created) and a proper package name: `mockgen -package mockdb -destination db/mock/store.go github.com/blessedmadukoma/go-simple-bank/db/sqlc Store` . Before running the command, get the mock package: `go get github.com/golang/mock`, remove the `// indirect` in go.mod
    4.  add mock command to Makefile
    5.  implement test for GetAccount API, added table-driven tests by testing multiple values in a struct (with fields `name` to identify the name of the test, `accountID`, `buildStubs`, and `checkResponse`)
    6.  to have cleaner logs for the tests output, we create a new `main_test.go` in the api folder and change the gin.SetMode to TestMode

14. Lecture 14: Custom params validator in Go
    1.  created transfer.go to handle transfer API
    2.  wrote a validAccount method to check if an account exists and the currency matches the input currency
    3.  updated server.go by adding the route for transfer
    4.  created a validator.go file to validate the currency, which removes the duplication we have
    5.  created currency.go in util package, added supported currency as constants and implemented a method `IsSupportedCurrency` to validate if the currency is supported or not.
    6.  register the custom validation method `validCurrency` in server.go
    7.  replace all oneOf... with the name of the registered validation: `currency` 

15. Lecture 15: Add users table with unique foreign key contraints in PostgreSQL
    1.  updated dbdiagram design by adding a user's table, linking the `owner` field in Accounts to the user's table and setting a composite index on the currency and owner meaning a user must only have one USD account (not 2) or one EUR account (not 2 EUR accts), but can have USD, EUR, NGN accounts.
    2.  generated a new migration add_users
    3.  replaced the composite index with: ALTER TABLE "accounts" ADD CONSTRAINT "owner_currency_key" UNIQUE ("owner", "currency") . They do the same thing, so pick one
    4.  updated Makefile by adding migrateup1 and migratedown1 commands to either migrate the very last or latest schema sql up or down respectively

16. Lecture 16: Handle DB errors in Go
    1.  create a user.sql file in db/query and write sql to create and get user by username. Generate sqlc.
    2.  Write tests to check the create and get user methods
    3.  update the failed tests such as accounts_test.go and transfers_test.go by adding the randomCreateUser method to the test files and replacing all `Owner` params with `user.Username` 
    4.  run `make mock` to regenerate the mockgen code to avoid errors
    5.  update account.go API to return a proper PostgreSQL error, if error found when creating an account

17. Lecture 17: Hash passwords in Go using Bcrypt
    1.  create a password.go and password_test.go in util package to handle Hashing of password, checking the hashed password with user input password, and testing both methods respectively.
    2.  update createRandomUser by adding the hashPassword method which generates a hashed password
    3.  created user.go API to create a new user and get a user by username

18. Lecture 18: Write stronger Golang unit tests with a custom Go-Mock matcher
    1.  quickly added test case for users in the `api` folder
    2.  wrote tests to validate the users details

19. Lecture 19: Why PASETO is better than JWT
    1.  Problems with JWT
        1.  weak algorithms: too much algorithms to choose from and some algorithms are known to be vulnerable
        2.  trivial forgery
    2. Why PASETO - Platform-Agnostic SEcurity TOkens
       1. Stronger algorithms - devs don't have to choose the algorithm, only select the version
       2. Non-trivial forgery - everything is authenticated, no more "none" algorithm
       3. PASETO structure:
          1. Local
            - main parts of the token: 
              - version
              - purpose e.g. local
              - payload:
                 - body (encrypted i.e. hashed or decrypted i.e. json format)
                 - Nonce
                 - Authentication tag
              - Footer (encoded or decoded)
          2. Public
           
           - main parts of the token: 
              - version
              - purpose e.g. public
              - payload:
                 - body (encrypted i.e. hashed or decrypted i.e. json format)
               - Signature

20. Lecture 20: How to create and verify JWT and PASETO token in Go
    1. created `token` package/folder to make the token
    2. wrote JWT token:
       1. created `Payload` struct to contain the data of the token, wrote `NewTokenPayload` function to create a new payload
       2. created `jwt` file to create JWT and the token, and verify the token
       3. wrote tests for JWTMaker (creating and validating token), expired token and invalid token also with no algorithm used
    3. Wrote Paseto token:
       1. created `paseto` file for making a new paseto token, creating paseto token and verifying the token
       2. wrote tests for `TestPasetoMaker`, `TestInvalidKeySize`, `TestExpiredPasetoToken` and `TestInvalidPasetoToken`


21. Lecture 21: Implement login API with PASETO and JWT
    1. added `tokenMaker` field to Server struct
    2. added new fields: `TokenSymmetricKey` and `AccessTokenDuration` to the `config` struct
    3. updated `NewServer` to take in config as a parameter, create a new token (PASETO)
    4. wrote a `newTestServer` for testing the `NewServer` method
    5. updated `api/user.go` by creating a `userResponse` struct and `login` method
    6. created a function `newUserResponse` which returns a newly created user
    7. wrote the `loginUser` handler api method to log in a user
    8. added  `/login` endpoint or route to the server
    9. grouped the routes
    10. either JWT or PASETO can be used to generate access token (comment out any of the lines in `api/server.go`)

22. How to implement authentication middleware
    1.  created `api/middleware.go`
    2.  wrote authentication middleware to verify if the access token is present in the authorization header
    3.  wrote test for authMiddleware function
    4.  updated specific routes to have authorization
    5.  removed `owner` field in `CreateAccountRequest` struct because of the auth route
    6.  added `authPayload` which contains the owner info from the `Authorization` token
    7.  updated `getAccountByID` to validate if the account gotten belongs to the unauthenticated user
    8.  updated `accounts.sql` by adding a WHERE clause to list accounts based on the authenticated user
    9.  ran `make sqlc` and `make mock` on the terminal to regenerate/update the database and mockgen i.e. `accounts.sql.go` file
    10. update `db/sqlc/TestListAccounts` test function by adding an owner to the account, the listAccountsParams and requiring the account should not be empty.
    11. updated `api/listAccounts` by adding an owner which is generated from the bearer token.
    12. modified transfer.go by updating the `validAccount` method through returning `account`, updating `createTransfer` method to check if the user is authenticated before making a transfer
    13. updated `account_test.go` by adding a new parameter `owner` to specify the authenticated owner, updated `TestGetAccountAPI` by modifying the create account function to generate a new user, added `setupAuth` function 
    14. fixed tests errors and added new tests...probably the craziest code I've written

23. How to build a minimal Golang Docker image
    1.  checked out to a new branch for building a docker image
    2.  updated golang version in `go.mod` and github workflow
    3.  created a Dockerfile containing:
         ```
         // use a lightweight Golang docker image <br>
         FROM golang:1.19.5-alpine3.16 <br>
         // set the current working directory to `app` <br>
         WORKDIR /app <br>
         // copy all the contents to the working directory <br>
         COPY . . <br>
         // build the contents with the executable being named `main`
         RUN go build -o main main.go <br>

         // expose the port 
         EXPOSE 9000 <br/>
         CMD [ "/app/main" ] <br/>
            
         ```
    4.  Built the Dockerfile using `docker build -t simplebank:latest .`
    5.  Due to the large ddocker image size generated, we converted the code to binary file:

      ```
         # Build stage
         FROM golang:1.19.5-alpine3.16 AS builder
         WORKDIR /app
         COPY . .
         RUN go build -o main main.go

         # Run stage
         FROM alpine:3.16
         WORKDIR /app
         COPY --from=builder /app/main .

         EXPOSE 9000
         CMD [ "/app/main" ]
      ```
        - and ran `docker build -t simplebank:latest .` again.
        - checked the image size: `docker images`

24. How to connect containers in the same docker network
    1.  start a docker container: docker run --name `container_name` -p `port:port` `image_name:tag` e.g. `docker run --name simplebank -p 9000:9000 simplebank:latest`
    2.   added `COPY app.env .` to Dockerfile which copies our environmental variables which will only be used for test purposes.
    3.   before re-running the docker container, remove the existing image: docker rm `name of container_image` e.g. `docker rm simplebank`. Then `docker ps -a` to view all containers. Then `docker rmi image_id` to remove the docker image e.g. `docker rmi e4c8ec76fe6a`. Check if the image still exists: `docker images`
    4.   rebuild the docker image: `docker build -t simplebank:latest .`
    5.   run the docker image: `docker run --name simplebank -p 9000:9000 simplebank:latest`.
    6.   stopped the server and removed the docker image: `docker rm simplebank`
    7.   ran the server with GIN in production mode: `docker run --name simplebank -e GIN_MODE=release -p 9000:9000 simplebank:latest`. Does not show any logs at first, until a request is made.
    8.   Run: `docker container inspect image_name` to view detailed information about an image such as Network settings. E.g. `docker container inspect postgres14`, `docker container inspect simplebank`
    9.   stopped the `simplebank` container, removed it (`docker rm simplebank`) and started it again but with an env for the database to connect to the docker postgres database: `docker run --name simplebank -p 9000:9000 -e GIN_MODE=release -e DB_SOURCE="postgresql://postgres:postgres@172.17.0.2:5432/simplebank?sslmode=disable" simplebank:latest`
    10.  Best way to connect to a docker Postgres db: using user defined network
         1.   run `docker network ls` to view the running networks and IDs
         2.   view more information of the specific network: `docker network inspect network_name` e.g. `docker network inspect bridge`
         3.   create my own network: `docker network create network_name` e.g. `docker network create bank-network`
         4.   connect the created network: `docker network connect network_name container_name` e.g. `docker network connect bank-network postgres14`
         5.   check the network details and copy the IP address of the container: `docker container inspect container_name` e.g. `docker container inspect postgres14`
         6.   run the simple bank docker application but with a network tag to run the  postgres database on the same network: `docker run --name simplebank --network bank-network -p 9000:9000 -e GIN_MODE=release -e DB_SOURCE="postgresql://postgres:postgres@postgres14/simplebank?sslmode=disable" simplebank:latest`
    11. updated postgres command in `Makefile` to run on the created network `bank-network`

25. How to use Docker compose
    1.  Docker compose will help in starting multiple services at once and control their start-up orders (wait-for)
    2.  created a docker-compose.yaml file, added services, environments and ports.
    3.  ran `docker compose up`
    4.  added the command to install golang-migrate to the Dockerfile
    5.  added the command which copies the migration files (located in db/migration) into the Dockerfile.
    6.  added the command which installs cURL since the alpine image by default does not come with cURL.
    7.  changed the way the app starts by creating a `start.sh` file to run the database migration before running the app and gave it executable permission (chmod +x start.sh)
    8.  added entry point for running the app to the Dockerfile
    9.  ran `docker compose down` to remove the existing containers and networks, then ran `docker compose up`
    10. added start.sh to Dockerfile, added `depends_on` tag in docker-compose.yaml, added wait-for.sh to work with `depends_on` tag

26-36. Skipped lectures 26 to 36: involves AWS and Kubernetes

37. Manage sessions with Refresh Token