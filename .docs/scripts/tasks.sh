#!/bin/bash
# Snake Game Task Manager (DDD)
# Команды: create, list, status <N>, done <N>, progress, next, changelog, help

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TASKS_DIR="$(dirname "$SCRIPT_DIR")/tasks"
BACKLOG="$TASKS_DIR/backlog.md"
CHANGELOG="$(dirname "$SCRIPT_DIR")/changelog.md"

get_next_num() {
    local max=0
    for f in "$TASKS_DIR"/[0-9][0-9][0-9]-*.md; do
        [ -f "$f" ] || continue
        local n=$(basename "$f" | cut -d'-' -f1 | tr -d '0')
        [ "$n" -gt "$max" ] 2>/dev/null && max=$n
    done
    printf "%03d" $((max + 1))
}

get_field() {
    sed 's/\r$//' "$2" | grep "^- \*\*$1\*\*:" 2>/dev/null | sed 's/.*: //' || true
}

case "${1:-help}" in
create)
    title="${2:?Укажите название задачи}"
    priority="${3:-🟡 Средний}"
    num=$(get_next_num)
    slug=$(echo "$title" | tr '[:upper:]' '[:lower:]' | tr ' ' '-' | tr -cd '[:alnum:]-')
    file="$TASKS_DIR/${num}-${slug}.md"
    today=$(date +%Y-%m-%d)
    cat > "$file" <<EOF
# #${num} — ${title}

- **Статус**: ⬜ Не начато
- **Приоритет**: ${priority}
- **Автор**: @username
- **Создана**: ${today}
- **Завершена**: —
- **Блокирует**: нет

## Описание
TODO: опишите задачу.

## Контекст
TODO: какие файлы затронуты.

**Файлы:**
- \`path/to/file.go\`

## Чек-лист
- [ ] Изучить код
- [ ] Реализовать
- [ ] Тесты
- [ ] go vet ./...
- [ ] go test ./...
- [ ] Обновить документацию
- [ ] Changelog

## Критерии приёмки
- [ ] Работает
- [ ] Нет регрессий

## Заметки
TODO.
EOF
    echo "✅ Создана: .docs/tasks/${num}-${slug}.md"
    ;;
list)
    echo "📋 Задачи:"
    grep '|' "$BACKLOG" | grep '[0-9]' | while IFS='|' read -r _ col_num col_task col_status col_pri col_files; do
        num=$(echo "$col_num" | tr -d '[:space:]')
        task=$(echo "$col_task" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
        status=$(echo "$col_status" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
        pri=$(echo "$col_pri" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
        printf "  %s | %-35s | %s | %s\n" "$num" "$task" "$status" "$pri"
    done
    ;;
status)
    num="${2:?Укажите номер}"
    file=$(ls "$TASKS_DIR"/${num}-*.md 2>/dev/null | head -1)
    [ -z "$file" ] && { echo "❌ Задача #$num не найдена"; exit 1; }
    echo "📌 $(head -1 "$file")"
    echo "   Статус: $(get_field 'Статус' "$file")"
    echo "   Приоритет: $(get_field 'Приоритет' "$file")"
    grep "^- \[" "$file" 2>/dev/null | sed 's/\r$//' | sed 's/^/   /' || true
    ;;
done)
    num="${2:?Укажите номер}"
    file=$(ls "$TASKS_DIR"/${num}-*.md 2>/dev/null | head -1)
    [ -z "$file" ] && { echo "❌ Не найдена"; exit 1; }
    title=$(head -1 "$file" | sed 's/\r$//' | sed 's/^# #*[0-9]* — //')
    today=$(date +%Y-%m-%d)
    sed -i.bak "s/^- \*\*Статус\*\*: .*/- **Статус**: ✅ Готово/" "$file"
    sed -i.bak "s/^- \*\*Завершена\*\*: .*/- **Завершена**: ${today}/" "$file"
    rm -f "${file}.bak"
    echo "✅ #${num}: $title — готово!"
    ;;
progress)
    total=0; done_count=0
    for f in "$TASKS_DIR"/[0-9][0-9][0-9]-*.md; do
        [ -f "$f" ] || continue
        total=$((total+1))
        sed 's/\r$//' "$f" | grep -q "✅ Готово" && done_count=$((done_count+1))
    done
    pct=0; [ "$total" -gt 0 ] && pct=$((done_count*100/total))
    echo "📊 Прогресс: $done_count/$total ($pct%)"
    ;;
next)
    for f in "$TASKS_DIR"/[0-9][0-9][0-9]-*.md; do
        [ -f "$f" ] || continue
        status=$(get_field 'Статус' "$f")
        if [[ "$status" == *"Не начато"* ]] || [[ "$status" == *"В работе"* ]]; then
            echo "➡️ $(head -1 "$f" | sed 's/\r$//')"
            echo "   Статус: $status"
            break
        fi
    done
    ;;
changelog) cat "$CHANGELOG" ;;
help|*)
    echo "Использование: ./tasks.sh {create|list|status|done|progress|next|changelog}"
    ;;
esac

