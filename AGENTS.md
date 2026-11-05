# AGENTS.md — server

Правила для агента в репозитории сервера. Общие правила workspace — в
родительском `AGENTS.md`, если он доступен.

## Стек

- Go 1.26, Gin, gorp, MSSQL, Redis, GLPI, Active Directory.
- Платформа разработки: Windows + PowerShell.

## Команды

Из каталога `server/`:

```powershell
gofmt -w <изменённые .go>
go build ./...
go vet ./...
go test ./...
```

Быстрый цикл: `go test ./internal/usecase/...`.
Перед финальным ответом — полный набор. При необходимости гарантированного
перезапуска использовать `go test -count=1 ./...`.

## Структура

- `internal/entity` — структуры, JSON/db-теги.
- `internal/usecase` — бизнес-логика, интерфейсы `Repository`/`UseCase`.
- `internal/adapter/repository/mssql` — MSSQL.
- `internal/adapter/repository/glpi` — GLPI.
- `internal/adapter/repository/ad` — Active Directory.
- `internal/adapter/repository/redisclient` — Redis.
- `internal/controller/http/api` — handlers и `*_rg.go` с маршрутами.
- `internal/app` — DI, запуск, фоновые циклы.
- `migrations` — SQL-миграции.

## Правила

- Новый endpoint вести по цепочке:
  `entity → usecase interface → usecase → repository → handler → route → tests/docs`.
- Интерфейсы обновлять в `internal/usecase/usecase.go` и
  `internal/controller/http/api/handler.go`.
- Миграции: `migrations/YYYYMMDD_описание.sql`, idempotent-стиль
  (`IF OBJECT_ID(...) IS NULL`, `IF NOT EXISTS`).
- Не применять миграции к живой БД; только создавать SQL-файлы.
- Сохранять существующий стиль запросов и placeholder-ов (`?`, gorp).
- Комментарии в коде — по-русски, как в проекте.
- Не откатывать несвязанные изменения и не переписывать чужие правки.
- Не создавать коммиты и не делать push без явной просьбы.
- Не логировать секреты, токены и персональные данные.
- Фоновые циклы должны проверять `IsAppInitialized()`.

## Ограничения

- Часть документации может лежать в родительском `../docs/`; если её нет,
  ориентироваться на код и `README`.
- Исторический мёртвый/закомментированный код не удалять без проверки.
