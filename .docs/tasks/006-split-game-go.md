# #006 — Разделить game.go на модули (engine/render/input)

- **Статус**: ✅ Выполнено
- **Приоритет**: 🔥 Высокий
- **Автор**: @rekhin
- **Создана**: 2026-11-26
- **Завершена**: 2026-11-26
- **Блокирует**: нет

## Описание

Файл `game.go` слишком большой (~500 строк). Разделён на логические модули: engine (игровой цикл), render (отрисовка), input (обработка ввода).

**Результат:**
- `game/game.go` — ~20 строк (цвета, Layout)
- `game/engine.go` — ~220 строк (Game struct, update loop, game logic, effects, cacheGrid)
- `game/render.go` — ~180 строк (Draw, drawMenu, drawGame, drawOverlay, drawHUD)
- `game/input.go` — ~50 строк (processInputPause, processInputDirection, processInput)

## Чек-лист

- [x] Изучить структуру game.go
- [x] Определить границы модулей
- [x] Создать engine.go — игровой цикл, update
- [x] Создать render.go — все draw* функции
- [x] Создать input.go — обработка ввода
- [x] Обновить импорты и зависимости
- [x] go vet ./...
- [x] go test ./...

## Критерии приёмки

- [x] game.go стал <200 строк (теперь ~20)
- [x] Все тесты проходят (30/30 PASS)
- [x] Код не сломался (go vet clean)
- [x] Архитектура стала понятнее

## Заметки

Важно сохранить публичный интерфейс Game — он используется в main.go.

### Исправленные ошибки при рефакторинге
- `/randSource` → `randSource` (опечатка в engine.go)
- `borderColor := ...` внутри цикла → объявлен до цикла с `var borderColor color.RGBA`
- Добавлены недостающие импорты: `ebiten`, `inpututil`, `vector`, `color`, `math`