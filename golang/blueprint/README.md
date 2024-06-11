# blueprint
First app!!!!!!!!!!!!!!

## Build and run

```
goreleaser --clean --snapshot
docker run -it -v $(pwd)/config.yaml:/app/config.yaml -p 8080:8080 ghcr.io/synkube/app/blueprint:latest /app/blueprint --config config.yaml
```
