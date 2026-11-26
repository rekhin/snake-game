package game

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

const highScoreFileName = "highscore.txt"

// highScoreFilePath возвращает путь к файлу сохранения рекорда.
// Использует os.UserConfigDir() для кроссплатформенного хранения:
//   - Windows: %APPDATA%\snake-game\highscore.txt
//   - Linux: ~/.config/snake-game/highscore.txt
//   - macOS: ~/Library/Application Support/snake-game/highscore.txt
//
// В случае ошибки получения UserConfigDir fallback на executable dir.
func highScoreFilePath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		// fallback: рядом с исполняемым файлом
		log.Printf("[highscore] не удалось получить UserConfigDir: %v, используем fallback", err)
		return filepath.Join(filepath.Dir(executablePath()), highScoreFileName)
	}
	return filepath.Join(configDir, "snake-game", highScoreFileName)
}

// executablePath возвращает путь к исполняемому файлу (с кэшированием).
var _executablePath string = func() string {
	p, _ := os.Executable()
	return p
}()

func executablePath() string {
	return _executablePath
}

// ensureHighScoreDir создаёт директорию для highscore файла, если она не существует.
func ensureHighScoreDir(path string) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Printf("[highscore] не удалось создать директорию %s: %v", dir, err)
	}
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
	ensureHighScoreDir(path)

	f, err := os.Create(path)
	if err != nil {
		log.Printf("[highscore] не удалось записать файл %s: %v", path, err)
		return
	}
	defer f.Close()

	fmt.Fprintln(f, g.score)
}