# blueprint
First app!!!!!!!!!!!!!!

## Build and run

```
goreleaser --clean --snapshot
docker run -it -v $(pwd)/config.yaml:/app/config.yaml -p 8080:8080 ghcr.io/synkube/app/blueprint:latest /app/blueprint --config config.yaml
```

## Run indexer
```
go run ./main.go --config config/config.yaml
```

## Run server (GraphQL)
```
go run ./main.go --config config/config_server.yaml server
```

## Postgres
```
docker run --name my-postgres -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=mypassword -e POSTGRES_DB=postgres -p 5432:5432 -d postgres
```

Psql commands
```
\l # list databases
\c postgres # connect to database
\dt # list tables
\q # quit
\? # help
select * from public.users;
```

## Clickhouse
```
docker run -d --name clickhouse-server -p 8123:8123 -p 9000:9000 clickhouse/clickhouse-server
```

## GraphQL
### Code generation
```
go run github.com/99designs/gqlgen generate
github.com/99designs/gqlgen/codegen/config@v0.17.49 # if necessary
```
