set dotenv-load := true

env name:
    ln -sf .env.{{name}} .env

dbinit:
    go run cmd/service/main.go db init

cleango:
    go clean -modcache