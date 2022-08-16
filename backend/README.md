# Backend
Backend приложения

* Установка пакетов (из папки backend/app)
- `go install backend/internal/app backend/internal/config backend/pkg/logging` — установка пакетов,

* Сборка и запуск (из папки backend)
- `go run ./app/cmd/server/main.go` — запуск,
- `goose create add_some_column sql` — создание миграции,
- `goose up` — Применить миграции
- example: `goose postgres "user=postgres dbname=stack sslmode=disable" up`
---
