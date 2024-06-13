# blueprint
First app!!!!!!!!!!!!!!

## Build and run

```
goreleaser --clean --snapshot
docker run -it -v $(pwd)/config.yaml:/app/config.yaml -p 8080:8080 ghcr.io/synkube/app/blueprint:latest /app/blueprint --config config.yaml
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
