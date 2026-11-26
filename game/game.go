package game

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"image/color"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/rekhin/snake-game/entity"
)

// Цвета игры
var (
	colorDarkBg    = color.RGBA{30, 30, 46, 255}
	colorGridLight = color.RGBA{40, 40, 58, 255}
	colorGridDark  = color.RGBA{35, 35, 51, 255}
	colorSnakeHead = color.RGBA{102, 217, 160, 255}
	colorSnakeBody = color.RGBA{48, 169, 118, 255}
	colorFood      = color.RGBA{235, 87, 87, 255}
	colorText      = color.RGBA{230, 230, 230, 255}
	colorTitle     = color.RGBA{102, 217, 160, 255}
	colorHighlight = color.RGBA{245, 196, 63, 255}
	colorHUD       = color.RGBA{20, 20, 30, 255}
	colorOverlay   = color.RGBA{0, 0, 0, 180}
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
	gridCache   *ebiten.Image // Кэшированное игровое поле (фон + сетка)
	shakeOffset [2]float32    // Смещение тряски экрана
	shakeTime   time.Duration // Оставшееся время тряски
	pulsePhase  float64       // Фаза пульсации еды
	randSource  *rand.Rand    // Локальный источник случайных чисел для эффектов
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

// resetGame сбрасывает состояние игры к начальному
func (g *Game) resetGame() {
	g.snake = entity.NewSnake()
	g.food = entity.NewFood(g.snake.Body)
	g.score = 0
	g.moveTimer = time.Now()
	g.inputQueue = nil

	// Сброс эффектов
	g.shakeOffset = [2]float32{0, 0}
	g.shakeTime = 0
	g.pulsePhase = 0

	// Пересоздаём кэш сетки
	g.cacheGrid()
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
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.resetGame()
		g.state = StatePlaying
	}
}

// updatePlaying обрабатывает состояние игры
func (g *Game) updatePlaying() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.state = StatePaused
		return
	}

	if g.snake != nil {
		// Обработка ввода (добавление в очередь)
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
			g.inputQueue = append(g.inputQueue, entity.DirUp)
		} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
			g.inputQueue = append(g.inputQueue, entity.DirDown)
		} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
			g.inputQueue = append(g.inputQueue, entity.DirLeft)
		} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
			g.inputQueue = append(g.inputQueue, entity.DirRight)
		}

		// Ограничение очереди (максимум 2 ввода для отзывчивости без задержек)
		if len(g.inputQueue) > 2 {
			g.inputQueue = g.inputQueue[len(g.inputQueue)-2:]
		}

		// Применение одного направления из очереди за один тик
		if len(g.inputQueue) > 0 {
			nextDir := g.inputQueue[0]
			g.inputQueue = g.inputQueue[1:]
			g.snake.SetDirection(nextDir)
		}
	}

	// Обновляем эффекты
	g.updateEffects()

	now := time.Now()
	if now.Sub(g.moveTimer) >= g.getMoveInterval() {
		g.moveSnake()
		g.moveTimer = now
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
		g.shakeTime = 300 * time.Millisecond
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

// updateEffects обновляет визуальные эффекты (тряска, пульсация)
func (g *Game) updateEffects() {
	now := time.Now()

	// Обновляем тряску экрана
	if g.shakeTime > 0 {
		g.shakeTime -= now.Sub(g.moveTimer) // approximate decay
		if g.shakeTime <= 0 {
			g.shakeTime = 0
			g.shakeOffset = [2]float32{0, 0}
		} else {
			// Случайное смещение для тряски
			g.shakeOffset[0] = (g.randSource.Float32() - 0.5) * 6
			g.shakeOffset[1] = (g.randSource.Float32() - 0.5) * 6
		}
	}

	// Обновляем фазу пульсации еды
	if g.pulsePhase > 0 {
		g.pulsePhase += 0.15
		if g.pulsePhase > math.Pi*2 {
			g.pulsePhase = 0
		}
	}
}

// Draw отрисовывает игру
func (g *Game) Draw(screen *ebiten.Image) {
	switch g.state {
	case StateMenu:
		g.drawMenu(screen)
	case StatePlaying:
		g.drawGame(screen)
	case StatePaused:
		g.drawGame(screen)
		g.drawOverlay(screen, "ПАУЗА")
	case StateGameOver:
		g.drawGame(screen)
		g.drawOverlay(screen, fmt.Sprintf("ИГРА ОКОНЧЕНА\nСчёт: %d", g.score))
	}
}

// drawMenu отрисовывает главное меню
func (g *Game) drawMenu(screen *ebiten.Image) {
	w, h := screen.Bounds().Dx(), screen.Bounds().Dy()
	screen.Fill(colorDarkBg)

	// Рисуем заголовок по центру
	titleText := "🐍 ЗМЕЙКА"
	titleFace := g.getFontFace(48)
	titleW, _ := text.Measure(titleText, titleFace, 0)
	top := &text.DrawOptions{}
	top.GeoM.Translate(float64(w/2)-titleW/2, float64(h/3))
	top.ColorScale.ScaleWithColor(colorTitle)
	text.Draw(screen, titleText, titleFace, top)

	// Рисуем подсказку по центру
	hintText := "Нажмите ENTER для начала"
	hintFace := g.getFontFace(24)
	hintW, _ := text.Measure(hintText, hintFace, 0)
	top = &text.DrawOptions{}
	top.GeoM.Translate(float64(w/2)-hintW/2, float64(h/2))
	top.ColorScale.ScaleWithColor(colorText)
	text.Draw(screen, hintText, hintFace, top)
}

// drawGame отрисовывает игровое поле с эффектами
func (g *Game) drawGame(screen *ebiten.Image) {
	if g.snake == nil || g.food == nil {
		return
	}

	screen.Fill(colorDarkBg)

	// Рисуем кэшированный фон с эффектом тряски
	op := &ebiten.DrawImageOptions{}
	if g.shakeTime > 0 {
		op.GeoM.Translate(float64(g.shakeOffset[0]), float64(g.shakeOffset[1]))
	}
	screen.DrawImage(g.gridCache, op)

	head := g.snake.GetHead()

	// Рисуем тело змейки (без головы — фон уже содержит сетку)
	for i := 1; i < len(g.snake.Body); i++ {
		pos := g.snake.Body[i]
		dx := float32(pos.X*CellSize)
		dy := float32(pos.Y*CellSize + CellSize + 4)
		vector.FillRect(screen, dx, dy, float32(CellSize), float32(CellSize), colorSnakeBody, false)

		// Граница сегмента тела
		borderColor := color.RGBA{255, 255, 255, 20}
		vector.StrokeLine(screen, dx, dy, dx+float32(CellSize), dy, 1, borderColor, false)
		vector.StrokeLine(screen, dx, dy+float32(CellSize)-1, dx+float32(CellSize), dy+float32(CellSize)-1, 1, borderColor, false)
		vector.StrokeLine(screen, dx, dy, dx, dy+float32(CellSize), 1, borderColor, false)
		vector.StrokeLine(screen, dx+float32(CellSize)-1, dy, dx+float32(CellSize)-1, dy+float32(CellSize), 1, borderColor, false)
	}

	// Рисуем голову змейки
	headX := float32(head.X * CellSize)
	headY := float32(head.Y*CellSize + CellSize + 4)
	vector.FillRect(screen, headX, headY, float32(CellSize), float32(CellSize), colorSnakeHead, false)

	// Граница головы
	borderColor := color.RGBA{255, 255, 255, 40}
	vector.StrokeLine(screen, headX, headY, headX+float32(CellSize), headY, 1.5, borderColor, false)
	vector.StrokeLine(screen, headX, headY+float32(CellSize)-1, headX+float32(CellSize), headY+float32(CellSize)-1, 1.5, borderColor, false)
	vector.StrokeLine(screen, headX, headY, headX, headY+float32(CellSize), 1.5, borderColor, false)
	vector.StrokeLine(screen, headX+float32(CellSize)-1, headY, headX+float32(CellSize)-1, headY+float32(CellSize), 1.5, borderColor, false)

	// Рисуем еду с эффектом пульсации
	foodPos := g.food.Position
	foodCenterX := float32(foodPos.X*CellSize + CellSize/2)
	foodCenterY := float32(foodPos.Y*CellSize+CellSize/2 + CellSize + 4)

	var foodColor = colorFood
	if g.pulsePhase > 0 {
		// Пульсация: меняем яркость и размер
		scale := 1.0 + math.Sin(g.pulsePhase)*0.15
		radius := float32(CellSize/2-2) * float32(scale)

		// Свечение вокруг еды
		glowColor := color.RGBA{235, 87, 87, 80}
		vector.DrawFilledCircle(screen, foodCenterX, foodCenterY, radius+4, glowColor, false)

		// Основная еда
		foodColor = color.RGBA{
			uint8(235 * (0.8 + math.Sin(g.pulsePhase)*0.2)),
			uint8(87 * (0.8 + math.Sin(g.pulsePhase)*0.2)),
			uint8(87 * (0.8 + math.Sin(g.pulsePhase)*0.2)),
			255,
		}
		vector.DrawFilledCircle(screen, foodCenterX, foodCenterY, radius, foodColor, false)

		// Блик на еде
		highlight := color.RGBA{255, 180, 180, 120}
		vector.DrawFilledCircle(screen, foodCenterX-2, foodCenterY-2, radius*0.3, highlight, false)
	} else {
		// Стандартная еда — круглая форма
		radius := float32(CellSize/2 - 2)
		vector.DrawFilledCircle(screen, foodCenterX, foodCenterY, radius, colorFood, false)

		// Блик на еде
		highlight := color.RGBA{255, 180, 180, 120}
		vector.DrawFilledCircle(screen, foodCenterX-2, foodCenterY-2, radius*0.3, highlight, false)
	}

	g.drawHUD(screen)
}

// drawOverlay рисует полупрозрачный оверлей с текстом
func (g *Game) drawOverlay(screen *ebiten.Image, overlayText string) {
	w, h := screen.Bounds().Dx(), screen.Bounds().Dy()

	// Оптимизировано: рисуем один прямоугольник вместо перебора всех пикселей
	vector.FillRect(screen, 0, 0, float32(w), float32(h), colorOverlay, false)

	lines := strings.Split(overlayText, "\n")
	centerY := h / 2
	lineHeight := 48
	startY := centerY - (len(lines)-1)*lineHeight/2

	for i, line := range lines {
		face := g.getFontFace(30)
		lineW, _ := text.Measure(line, face, 0)
		top := &text.DrawOptions{}
		top.GeoM.Translate(float64(w/2)-lineW/2, float64(startY+i*lineHeight))
		top.ColorScale.ScaleWithColor(colorText)
		text.Draw(screen, line, face, top)
	}
}

// drawHUD отрисовывает интерфейс пользователя
func (g *Game) drawHUD(screen *ebiten.Image) {
	if g.snake == nil {
		return
	}

	w, _ := screen.Bounds().Dx(), screen.Bounds().Dy()
	hudHeight := CellSize + 4

	// Оптимизировано: рисуем один прямоугольник вместо перебора всех пикселей
	vector.FillRect(screen, 0, 0, float32(w), float32(hudHeight), colorHUD, false)

	const padding = 8
	const rightPadding = 8

	top := &text.DrawOptions{}
	top.ColorScale.ScaleWithColor(colorText)
	text.Draw(screen, fmt.Sprintf("Счёт: %d", g.score), g.getFontFace(18), top)

	if g.highScore > 0 {
		top = &text.DrawOptions{}
		highScoreText := fmt.Sprintf("Рекорд: %d", g.highScore)
		highScoreFace := g.getFontFace(18)
		top.GeoM.Translate(float64(padding+90), 16)
		top.ColorScale.ScaleWithColor(colorHighlight)
		text.Draw(screen, highScoreText, highScoreFace, top)
	}

	length := len(g.snake.Body)
	lengthText := fmt.Sprintf("Длина: %d", length)
	lengthFace := g.getFontFace(18)
	lengthW, _ := text.Measure(lengthText, lengthFace, 0)

	top = &text.DrawOptions{}
	top.GeoM.Translate(float64(w)-lengthW-rightPadding, 16)
	top.ColorScale.ScaleWithColor(colorText)
	text.Draw(screen, lengthText, lengthFace, top)
}

// Layout возвращает размер игрового экрана
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return GridWidth * CellSize, GridHeight*CellSize + CellSize + 4
}
