set dotenv-load := true

default:
  @just --list

env name:
    ln -sf .env.{{name}} .env
    
dbinit:
    go run cmd/service/main.go db init
    
new_migration name:
    go run cmd/service/main.go db create-sql {{name}}

migrate:
    go run cmd/service/main.go db migrate

rollback:
    go run cmd/service/main.go db rollback

cleango:
    go clean -modcache