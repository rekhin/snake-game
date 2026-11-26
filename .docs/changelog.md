# 📝 Changelog — Snake Game

## [Unreleased]

### ✅ Завершено
- #002 — Перенести highscore в `os.UserConfigDir()` (2026-11-26)
  - `highScoreFilePath()` теперь использует `os.UserConfigDir()` для кроссплатформенного хранения
  - На Windows: `%APPDATA%\snake-game\highscore.txt`
  - На Linux: `~/.config/snake-game/highscore.txt`
  - На macOS: `~/Library/Application Support/snake-game/highscore.txt`
  - Fallback на директорию исполняемого файла при ошибке получения UserConfigDir
  - Добавлена функция `ensureHighScoreDir()` с `os.MkdirAll` для создания директории
  - Добавлено логирование ошибок (log.Printf) вместо silent игнорирования
  - Добавлены тесты: `TestHighScoreFilePath`, `TestEnsureHighScoreDir`
  - Все 30 тестов проходят (`go test ./...`)
- #001 — Заменить `math/rand` на локальный источник (2026-11-26)
  - Создан локальный `randSource *rand.Rand` в `entity/food.go` с инициализацией через `init()`
  - Добавлено поле `randSource` в структуру `Game` в `game/game.go`, инициализируется в `NewGame()`
  - Заменены все вызовы `rand.Intn()` на `randSource.Intn()` в `entity/food.go`
  - Заменены все вызовы `rand.Float32()` на `g.randSource.Float32()` в `game/game.go`
  - Добавлена функция `SetRandSource()` для инъекции детерминированного источника в тестах
  - Написаны 3 новых теста на детерминизм спавна еды
  - Исправлены существующие тесты `game/` для учёта нового поля `randSource`
  - Все 28 тестов проходят (`go test ./...`)

### ⬜ Запланировано
- #003 — Удалить мёртвый код `drawCell`
- #004 — Добавить nil-guard в `drawHUD`
- #005 — Исправить decay тряски на `time.Since`
- #006 — Разделить `game.go` на модули
- #007 — Добавить fallback для шрифта
- #008 — Улучшить тесты `game/`
- #009 — Добавить CI/CD
- #010 — Добавить конфигурацию через файл

---

## [v0.1.0] — 2024-11-23

### 🎉 Первая версия
- Классическая игра "Змейка" на Go + Ebitengine
- Размер поля: 30×25 клеток
- Управление: WASD + стрелки
- Состояния: меню, игра, пауза, game over
- Очки и рекорд (сохранение в файл)
- Визуальные эффекты: тряска экрана, пульсация еды
- Покрытие тестами: entity/ — 20 тестов, game/ — 8 тестов

---

## 📋 История изменений по задачам

> Задачи будут добавляться сюда по мере их завершения. Формат:
> 
> ### ✅ #NNN — Название задачи (YYYY-MM-DD)
> - Что было сделано
> - Какие файлы изменены
> - Ссылка на PR/коммит
