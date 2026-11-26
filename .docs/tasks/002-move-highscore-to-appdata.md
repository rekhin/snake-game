# #002 — Перенести highscore в `os.UserConfigDir()`

- **Статус**: ⬜ Не начато
- **Приоритет**: 🔥 Высокий
- **Автор**: @rekhin
- **Создана**: 2026-11-26
- **Завершена**: —
- **Блокирует**: нет

## Описание

Файл highscore.txt сохраняется рядом с исполняемым файлом (`filepath.Dir(executable)`). На Windows это может быть защищённая директория (Program Files), и запись не получится. Нужно перенести в `os.UserConfigDir()`.

## Контекст

Обнаружено при ревью проекта. Текущая реализация:
- `game/highscore.go:12-18` — `highScoreFilePath()` возвращает путь рядом с .exe
- Нет fallback в `%APPDATA%` или `%LOCALAPPDATA%` на Windows
- На Linux — `$HOME/.config`, на macOS — `~/Library/Application Support`

**Файлы:**
- `game/highscore.go` — строки 12-18, 47: функция `highScoreFilePath()`

## Чек-лист

- [ ] Изучить текущую реализацию `highScoreFilePath()` в highscore.go
- [ ] Заменить на использование `os.UserConfigDir()` с fallback на executable dir
- [ ] Использовать кроссплатформенный путь: `{UserConfigDir}/snake-game/highscore.txt`
- [ ] Создать директорию, если не существует (`os.MkdirAll`)
- [ ] Добавить логирование ошибки создания директории (но не панику)
- [ ] Обновить тест в `game_test.go` — проверить путь сохранения
- [ ] Проверка: `go vet ./...` проходит без ошибок
- [ ] Проверка: `go test ./...` все тесты зелёные
- [ ] Обновить документацию (README, комментарии)
- [ ] Записать в changelog.md

## Критерии приёмки

- [ ] Файл highscore сохраняется в UserConfigDir (не рядом с .exe)
- [ ] На Windows — `%APPDATA%\snake-game\highscore.txt`
- [ ] На Linux — `~/.config/snake-game/highscore.txt`
- [ ] На macOS — `~/Library/Application Support/snake-game/highscore.txt`
- [ ] Фоллбэк на executable dir при ошибке UserConfigDir
- [ ] Нет ошибок записи, нет паник
- [ ] Все тесты проходят

## Заметки

- `os.UserConfigDir()` доступен с Go 1.0
- Нужно убедиться, что директория создаётся перед записью
- Можно добавить флаг `-highscore-path` для отладки
