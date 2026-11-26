# #002 — Перенести highscore в `os.UserConfigDir()`

- **Статус**: ✅ Завершено (2026-11-26)
- **Приоритет**: 🔥 Высокий
- **Автор**: @rekhin
- **Создана**: 2026-11-26
- **Завершена**: 2026-11-26
- **Блокирует**: нет

## Описание

Файл highscore.txt сохранялся рядом с исполняемым файлом (`filepath.Dir(executable)`). На Windows это может быть защищённая директория (Program Files), и запись не получалась. Перенесено в `os.UserConfigDir()`.

## Контекст

Обнаружено при ревью проекта. Текущая реализация:
- `game/highscore.go:12-18` — `highScoreFilePath()` возвращает путь рядом с .exe
- Нет fallback в `%APPDATA%` или `%LOCALAPPDATA%` на Windows
- На Linux — `$HOME/.config`, на macOS — `~/Library/Application Support`

**Файлы:**
- `game/highscore.go` — строки 12-18, 47: функция `highScoreFilePath()`

## Чек-лист

- [x] Изучить текущую реализацию `highScoreFilePath()` в highscore.go
- [x] Заменить на использование `os.UserConfigDir()` с fallback на executable dir
- [x] Использовать кроссплатформенный путь: `{UserConfigDir}/snake-game/highscore.txt`
- [x] Создать директорию, если не существует (`os.MkdirAll`)
- [x] Добавить логирование ошибки создания директории (но не панику)
- [x] Обновить тест в `game_test.go` — проверить путь сохранения
- [x] Проверка: `go vet ./...` проходит без ошибок
- [x] Проверка: `go test ./...` все тесты зелёные
- [x] Обновить документацию (README, комментарии)
- [x] Записать в changelog.md

## Критерии приёмки

- [x] Файл highscore сохраняется в UserConfigDir (не рядом с .exe)
- [x] На Windows — `%APPDATA%\snake-game\highscore.txt`
- [x] На Linux — `~/.config/snake-game/highscore.txt`
- [x] На macOS — `~/Library/Application Support/snake-game/highscore.txt`
- [x] Фоллбэк на executable dir при ошибке UserConfigDir
- [x] Нет ошибок записи, нет паник
- [x] Все тесты проходят

## Заметки

- `os.UserConfigDir()` доступен с Go 1.0
- Нужно убедиться, что директория создаётся перед записью
- Можно добавить флаг `-highscore-path` для отладки
