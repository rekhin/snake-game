# #005 — Исправить decay тряски на `time.Since`

- **Статус**: ✅ Готово
- **Приоритет**: 💤 Низкий
- **Автор**: @rekhin
- **Создана**: 2026-11-26
- **Завершена**: 2026-11-26
- **Блокирует**: нет

## Описание

Эффект тряски (shake) использует некорректный decay. Нужно заменить на `time.Since` для плавного затухания.

**Файлы:**
- `game/game.go` — эффекты тряски

## Чек-лист

- [x] Изучить проблему с `shakeTime` и `now.Sub(g.moveTimer)`
- [x] Заменить поле `shakeTime time.Duration` на `shakeEndTime time.Time` (абсолютное время)
- [x] Обновить `moveSnake`: устанавливать `g.shakeEndTime = time.Now().Add(300ms)` вместо `g.shakeTime = 300ms`
- [x] Обновить `updateEffects`: использовать `time.Since` / `now.After(g.shakeEndTime)` для корректного decay
- [x] Обновить `resetGame`: инициализировать `shakeEndTime = time.Time{}`
- [x] Обновить `drawGame`: заменить проверку `g.shakeTime > 0` на `!g.shakeEndTime.IsZero()`
- [x] Обновить тесты: переписать `TestUpdateEffects_Shake` под новую логику
- [x] Проверка: `go vet ./...` проходит без ошибок
- [x] Проверка: `go test ./...` все тесты зелёные
- [x] Обновить changelog.md

## Критерии приёмки

- [x] `shakeTime` заменён на `shakeEndTime time.Time`
- [x] Тряска экрана корректно длится 300ms независимо от FPS/частоты вызова Update()
- [x] Нет паник и гонок данных
- [x] Проект компилируется без ошибок
- [x] Все тесты проходят
- [x] Пульсация еды тоже использует `time.Since` для корректного decay

## Заметки

Аналогично исправлена пульсация (`pulsePhase`) — теперь она считает фазу через `elapsed.Seconds() * 3` от `pulseStartTime`.
