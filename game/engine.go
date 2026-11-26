package game

import (
	"bytes"

	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/rekhin/snake-game/entity"
)

// Game представляет основную игровую структуру
type Game struct {
	state     State
	snake     *entity.Snake
	food      *entity.Food
	score     int
	highScore int
	moveTimer time.Time
	fontFace  text.Face
	inputQueue []entity.Direction

	// Эффекты
	gridCache      *ebiten.Image // Кэшированное игровое поле (фон + сетка)
	shakeOffset    [2]float32    // Смещение тряски экрана
	shakeEndTime   time.Time     // Абсолютное время окончания тряски
	pulsePhase     float64       // Фаза пульсации еды
	pulseStartTime time.Time     // Время начала пульсации
randSource    *rand.Rand    // Локальный источник случайных чисел для эффектов
}

// NewGame создаёт новую игру с заданным шрифтом (TTF байты)
func NewGame(fontData []byte) *Game {
	faceSource, err := text.NewGoTextFaceSource(bytes.NewReader(fontData))
	if err != nil {
		panic(fmt.Sprintf("failed to create font face: %v", err))
	}

	return &Game{
		state:     StateMenu,
		highScore: loadHighScore(),
		fontFace:  &text.GoTextFace{Source: faceSource, Size: 16},
		randSource: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// getFontFace возвращает Face с заданным размером
func (g *Game) getFontFace(size float64) text.Face {
	return &text.GoTextFace{Source: g.fontFace.(*text.GoTextFace).Source, Size: size}
}

// Update выполняет игровой цикл
func (g *Game) Update() error {
	switch g.state {
	case StateMenu:
		g.updateMenu()
	case StatePlaying:
		g.updatePlaying()
	case StatePaused:
		g.updatePaused()
	case StateGameOver:
		g.updateGameOver()
	}
	return nil
}

// updateMenu обрабатывает состояние меню
func (g *Game) updateMenu() {
	if g.processInput() {
		return
	}
	g.resetGame()
	g.state = StatePlaying
}

// updatePlaying обрабатывает состояние игры
func (g *Game) updatePlaying() {
	if g.processInputPause() {
		return
	}

	if g.snake != nil {
		g.processInputDirection()
	}

	// Обновляем эффекты
	g.updateEffects()

	now := time.Now()
	if now.Sub(g.moveTimer) >= g.getMoveInterval() {
		g.moveSnake()
		g.moveTimer = now
	}
}

// updatePaused обрабатывает состояние паузы
func (g *Game) updatePaused() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.state = StatePlaying
	}
}

// updateGameOver обрабатывает состояние game over
func (g *Game) updateGameOver() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.state = StateMenu
	}
}

// getMoveInterval возвращает интервал между ходами (зависит от счёта)
func (g *Game) getMoveInterval() time.Duration {
	interval := StartMoveIntervalMs - (g.score/SpeedUpEveryNPoints)*10
	if interval < MinMoveIntervalMs {
		interval = MinMoveIntervalMs
	}
	return time.Duration(interval) * time.Millisecond
}

// moveSnake выполняет ход змейки
func (g *Game) moveSnake() {
	if g.snake == nil || g.food == nil {
		return
	}

	g.snake.Move()

	if g.snake.CheckWallCollision() || g.snake.CheckSelfCollision() {
		g.state = StateGameOver
		// Тряска экрана при проигрыше
		g.shakeEndTime = time.Now().Add(300 * time.Millisecond)
		if g.score > g.highScore {
			g.highScore = g.score
			g.saveHighScore()
		}
		return
	}

	head := g.snake.GetHead()
	if head.Equals(g.food.Position) {
		g.snake.GrowBy(1)
		g.score += FoodScore
		g.food.Spawn(g.snake.Body)
		// Пульсация при поедании еды
		g.pulsePhase = 0
		g.pulseStartTime = time.Now()
	}
}

// updateEffects обновляет визуальные эффекты (тряска, пульсация)
func (g *Game) updateEffects() {
	now := time.Now()

	// Обновляем тряску экрана по абсолютному времени окончания
	if !g.shakeEndTime.IsZero() && now.After(g.shakeEndTime) {
		g.shakeEndTime = time.Time{}
		g.shakeOffset = [2]float32{0, 0}
	} else if !g.shakeEndTime.IsZero() {
		// Случайное смещение для тряски
		g.shakeOffset[0] = (g.randSource.Float32() - 0.5) * 6
		g.shakeOffset[1] = (g.randSource.Float32() - 0.5) * 6
	}

	// Обновляем фазу пульсации еды с time.Since для корректного decay
	if g.pulsePhase > 0 {
		elapsed := now.Sub(g.pulseStartTime)
		g.pulsePhase = elapsed.Seconds() * 3 // ~2 секунды на полный цикл
		if g.pulsePhase > math.Pi*2 {
			g.pulsePhase = 0
		}
	}
}

// resetGame сбрасывает состояние игры к начальному
func (g *Game) resetGame() {
	g.snake = entity.NewSnake()
	g.food = entity.NewFood(g.snake.Body)
	g.score = 0
	g.moveTimer = time.Now()
	g.inputQueue = nil

	// Сброс эффектов
	g.shakeOffset = [2]float32{0, 0}
	g.shakeEndTime = time.Time{}
	g.pulsePhase = 0
	g.pulseStartTime = time.Time{}

	// Пересоздаём кэш сетки
	g.cacheGrid()
}

// cacheGrid создаёт кэшированное изображение статического фона (сетка + HUD)
func (g *Game) cacheGrid() {
	w := GridWidth * CellSize
	h := GridHeight * CellSize
	cache := ebiten.NewImage(w, h)

	for x := 0; x < GridWidth; x++ {
		for y := 0; y < GridHeight; y++ {
			var cellColor color.RGBA
			if (x+y)%2 == 0 {
				cellColor = colorGridLight
			} else {
				cellColor = colorGridDark
			}

			dx := float32(x * CellSize)
			dy := float32(y * CellSize)
			vector.FillRect(cache, dx, dy, float32(CellSize), float32(CellSize), cellColor, false)

			// Границы клеток
			borderColor := color.RGBA{255, 255, 255, 30}
			vector.StrokeLine(cache, dx, dy, dx+float32(CellSize), dy, 1, borderColor, false)
			vector.StrokeLine(cache, dx, dy+float32(CellSize)-1, dx+float32(CellSize), dy+float32(CellSize)-1, 1, borderColor, false)
			vector.StrokeLine(cache, dx, dy, dx, dy+float32(CellSize), 1, borderColor, false)
			vector.StrokeLine(cache, dx+float32(CellSize)-1, dy, dx+float32(CellSize)-1, dy+float32(CellSize), 1, borderColor, false)
		}
	}

	g.gridCache = cache
}
