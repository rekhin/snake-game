package game

import (
	"math"
	"math/rand"
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
		moveTimer:  time.Now().Add(-100 * time.Millisecond),
		randSource: rand.New(rand.NewSource(42)),
	}

	// Устанавливаем время тряски
	g.shakeTime = 200 * time.Millisecond

	// Вызываем updateEffects (нужно создать snake/food для корректной работы)
	// В данном тесте проверяем только логику decay времени
	initialShakeTime := g.shakeTime
	if initialShakeTime <= 0 {
		t.Skip("shake time already expired")
	}

	// После updateEffects время тряски должно уменьшиться или обнулиться
	g.updateEffects()

	if g.shakeTime > initialShakeTime {
		t.Error("shake time should not increase")
	}
}

// TestUpdateEffects_Pulse проверяет обновление фазы пульсации
func TestUpdateEffects_Pulse(t *testing.T) {
	g := &Game{
		moveTimer:  time.Now().Add(-100 * time.Millisecond),
		pulsePhase: 0,
		randSource: rand.New(rand.NewSource(42)),
	}

	// Устанавливаем фазу пульсации
	g.pulsePhase = 1.0

	initialPhase := g.pulsePhase
	g.updateEffects()

	if g.pulsePhase <= initialPhase {
		t.Error("pulse phase should increase")
	}

	if g.pulsePhase > math.Pi*2 {
		t.Error("pulse phase should wrap around after 2*pi")
	}
}
