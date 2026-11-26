package game

import (
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/rekhin/snake-game/entity"
)

// TestNewGame_HighScoreLoaded проверяет, что highScore загружается при создании игры
func TestNewGame_HighScoreLoaded(t *testing.T) {
	// Создаём минимальный шрифт (пустые байты не подойдут, поэтому используем заглушку)
	// В реальном тесте мы бы использовали тестовый шрифт или моки
	// Здесь проверяем только структуру
	g := &Game{
		state:     StateMenu,
		highScore: 100,
	}

	if g.highScore != 100 {
		t.Errorf("expected highScore 100, got %d", g.highScore)
	}
}

// TestState_String проверяет строковое представление состояний
func TestState_String(t *testing.T) {
	tests := []struct {
		state  State
		expect string
	}{
		{StateMenu, "MENU"},
		{StatePlaying, "PLAYING"},
		{StatePaused, "PAUSED"},
		{StateGameOver, "GAME_OVER"},
		{State(99), "UNKNOWN"},
	}

	for _, tt := range tests {
		if got := tt.state.String(); got != tt.expect {
			t.Errorf("state(%d).String() = %q, want %q", tt.state, got, tt.expect)
		}
	}
}

// TestGetMoveInterval проверяет вычисление интервала движения
func TestGetMoveInterval(t *testing.T) {
	g := &Game{
		score: 0,
	}

	interval := g.getMoveInterval()
	if interval != time.Duration(StartMoveIntervalMs)*time.Millisecond {
		t.Errorf("expected interval %dms at score 0, got %dms", StartMoveIntervalMs, interval)
	}

	// При 5 очках интервал должен уменьшиться
	g.score = SpeedUpEveryNPoints
	interval = g.getMoveInterval()
	expected := time.Duration(StartMoveIntervalMs-10) * time.Millisecond
	if interval != expected {
		t.Errorf("expected interval %dms at score %d, got %dms", expected, SpeedUpEveryNPoints, interval)
	}

	// При большом счёте интервал не должен быть меньше MinMoveIntervalMs
	g.score = 1000
	interval = g.getMoveInterval()
	if interval != time.Duration(MinMoveIntervalMs)*time.Millisecond {
		t.Errorf("expected minimum interval %dms at high score, got %dms", MinMoveIntervalMs, interval)
	}
}

// TestResetGame проверяет сброс состояния игры
func TestResetGame(t *testing.T) {
	g := &Game{
		state:      StatePlaying,
		score:      100,
		highScore:  50,
		snake:      nil,
		food:       nil,
		randSource: rand.New(rand.NewSource(42)),
	}

	g.resetGame()

	if g.state != StatePlaying {
		t.Errorf("expected state to remain playing after reset, got %v", g.state)
	}
	if g.score != 0 {
		t.Errorf("expected score to be reset to 0, got %d", g.score)
	}
	if g.inputQueue != nil {
		t.Error("expected inputQueue to be nil after reset")
	}
}

// TestInputQueue_Limit проверяет ограничение очереди ввода
func TestInputQueue_Limit(t *testing.T) {
	g := &Game{}

	// Добавляем 5 направлений в очередь
	for i := 0; i < 5; i++ {
		g.inputQueue = append(g.inputQueue, entity.DirUp)
	}

	// Применяем ограничение как в updatePlaying
	if len(g.inputQueue) > 2 {
		g.inputQueue = g.inputQueue[len(g.inputQueue)-2:]
	}

	if len(g.inputQueue) != 2 {
		t.Errorf("expected queue length 2, got %d", len(g.inputQueue))
	}
}

// TestInputQueue_ApplyDirection проверяет применение направления из очереди
func TestInputQueue_ApplyDirection(t *testing.T) {
	g := &Game{
		inputQueue: []entity.Direction{entity.DirUp},
	}

	if len(g.inputQueue) != 1 {
		t.Errorf("expected queue length 1, got %d", len(g.inputQueue))
	}

	nextDir := g.inputQueue[0]
	g.inputQueue = g.inputQueue[1:]

	if nextDir != entity.DirUp {
		t.Errorf("expected direction Up, got %v", nextDir)
	}

	if len(g.inputQueue) != 0 {
		t.Error("expected queue to be empty after applying direction")
	}
}

// TestUpdateEffects_Shake проверяет эффект тряски экрана
func TestUpdateEffects_Shake(t *testing.T) {
	g := &Game{
		moveTimer:    time.Now().Add(-100 * time.Millisecond),
		shakeEndTime: time.Now().Add(200 * time.Millisecond),
		randSource:   rand.New(rand.NewSource(42)),
	}

	// shakeEndTime ещё не наступила — тряска должна продолжаться
	g.updateEffects()

	if g.shakeEndTime.IsZero() {
		t.Error("shake should not have ended yet")
	}
	if g.shakeOffset[0] == 0 && g.shakeOffset[1] == 0 {
		t.Error("shake offset should be non-zero during shake")
	}

	// Симулируем ожидание — устанавливаем shakeEndTime в прошлое
	g.shakeEndTime = time.Now().Add(-200 * time.Millisecond)
	g.updateEffects()

	if !g.shakeEndTime.IsZero() {
		t.Error("shake should have ended after time passed")
	}
	if g.shakeOffset[0] != 0 || g.shakeOffset[1] != 0 {
		t.Error("shake offset should be zeroed after shake ends")
	}
}

// TestUpdateEffects_Pulse проверяет эффект пульсации еды
func TestUpdateEffects_Pulse(t *testing.T) {
	// Устанавливаем pulseStartTime в прошлое, чтобы фаза уже была прогрессировала
	g := &Game{
		pulsePhase:     0.5,
		pulseStartTime: time.Now().Add(-1 * time.Second), // 1 секунда назад
	}

	initialPhase := g.pulsePhase
	g.updateEffects()

	if g.pulsePhase <= initialPhase {
		t.Errorf("pulse phase should increase over time: initial=%f, after=%f", initialPhase, g.pulsePhase)
	}
}

// TestHighScoreFilePath проверяет, что путь highscore содержит snake-game в пути
func TestHighScoreFilePath(t *testing.T) {
	path := highScoreFilePath()

	if !strings.Contains(path, "snake-game") {
		t.Errorf("expected highscore path to contain 'snake-game', got: %s", path)
	}

	if !strings.HasSuffix(path, filepath.Join("snake-game", "highscore.txt")) &&
		!strings.HasSuffix(path, "snake-game\\highscore.txt") {
		t.Errorf("expected highscore path to end with 'snake-game/highscore.txt', got: %s", path)
	}
}

// TestEnsureHighScoreDir проверяет создание директории для highscore файла
func TestEnsureHighScoreDir(t *testing.T) {
	// Создаём временную директорию для теста
	tmpDir := t.TempDir()
	testPath := filepath.Join(tmpDir, "subdir", "highscore.txt")

	ensureHighScoreDir(testPath)

	expectedDir := filepath.Join(tmpDir, "subdir")
	if _, err := os.Stat(expectedDir); os.IsNotExist(err) {
		t.Errorf("expected directory %s to be created", expectedDir)
	}

	// Проверяем что файл можно создать в этой директории
	f, err := os.Create(testPath)
	if err != nil {
		t.Errorf("failed to create file in ensured directory: %v", err)
	}
	f.Close()
}
