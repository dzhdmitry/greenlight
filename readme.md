# Requrements

- go 1.16 or higer

- migrate. install:

    ```bash
    $ curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz
    $ mv migrate.linux-amd64 $GOPATH/bin/migrate
    ```
    
- Postgresql database setup

    ```sql
    CREATE DATABASE greenlight;
    CREATE ROLE greenlight WITH LOGIN PASSWORD 'pa55word';
    CREATE EXTENSION IF NOT EXISTS citext;
    ```

- Migrate

    ```bash
    migrate -path=./migrations -database=postgres://greenlight:pa55word@localhost/greenlight?sslmode=disable up
    ```

# Run

go run ./cmd/api
