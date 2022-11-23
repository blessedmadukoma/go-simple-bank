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
   2. 
4. Lesson 4: