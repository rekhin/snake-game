package game

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

// highScoreFile возвращает путь к файлу сохранения рекорда
func highScoreFilePath() string {
	execPath, err := os.Executable()
	if err != nil {
		// fallback: текущая директория
		return "highscore.txt"
	}
	return filepath.Join(filepath.Dir(execPath), "highscore.txt")
}

// loadHighScore загружает сохранённый рекорд из файла
func loadHighScore() int {
	path := highScoreFilePath()
	f, err := os.Open(path)
	if err != nil {
		return 0 // файл не существует — возвращаем 0
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		score, err := strconv.Atoi(line)
		if err == nil && score > 0 {
			return score
		}
	}
	return 0
}

// saveHighScore сохраняет рекорд в файл
func (g *Game) saveHighScore() {
	if g.score <= g.highScore {
		return // новый рекорд не побит — не сохраняем
	}

	path := highScoreFilePath()
	f, err := os.Create(path)
	if err != nil {
		return // игнорируем ошибки записи
	}
	defer f.Close()

	fmt.Fprintln(f, g.score)
}