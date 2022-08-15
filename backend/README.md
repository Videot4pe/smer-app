# Backend
Backend приложения

- `go run ./app/cmd/server/main.go` — запуск,
- `goose create add_some_column sql` — создание миграции,
- `goose up` — Применить миграции
- example: `goose postgres "user=postgres dbname=stack sslmode=disable" up`
---
