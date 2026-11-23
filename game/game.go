package game

import (
	"bytes"
	"fmt"
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
}

// NewGame создаёт новую игру с заданным шрифтом (TTF байты)
func NewGame(fontData []byte) *Game {
	faceSource, err := text.NewGoTextFaceSource(bytes.NewReader(fontData))
	if err != nil {
		panic(fmt.Sprintf("failed to create font face: %v", err))
	}

	return &Game{
		state:    StateMenu,
		fontFace: &text.GoTextFace{Source: faceSource, Size: 16},
	}
}

// getFontFace возвращает Face с заданным размером
func (g *Game) getFontFace(size float64) text.Face {
	return &text.GoTextFace{Source: g.fontFace.(*text.GoTextFace).Source, Size: size}
}

// resetGame сбрасывает состояние игры к начальному
func (g *Game) resetGame() {
	g.snake = entity.NewSnake()
	g.food = entity.NewFood(g.snake.Body)
	g.score = 0
	g.moveTimer = time.Now()
	g.inputQueue = nil
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

		// Ограничение очереди, чтобы избежать задержек
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
		if g.score > g.highScore {
			g.highScore = g.score
		}
		return
	}

	head := g.snake.GetHead()
	if head.Equals(g.food.Position) {
		g.snake.GrowBy(1)
		g.score += FoodScore
		g.food.Spawn(g.snake.Body)
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

// drawGame отрисовывает игровое поле
func (g *Game) drawGame(screen *ebiten.Image) {
	screen.Fill(colorDarkBg)

	if g.snake == nil || g.food == nil {
		return
	}

	for x := 0; x < GridWidth; x++ {
		for y := 0; y < GridHeight; y++ {
			pos := entity.Point{X: x, Y: y}
			head := g.snake.GetHead()

			var cellColor color.RGBA
			switch {
			case head.Equals(pos):
				cellColor = colorSnakeHead
			case g.snake.IsBody(pos):
				cellColor = colorSnakeBody
			case pos.Equals(g.food.Position):
				cellColor = colorFood
			case (x+y)%2 == 0:
				cellColor = colorGridLight
			default:
				cellColor = colorGridDark
			}

			drawCell(screen, pos, cellColor)
		}
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
		top.GeoM.Translate(200, 16)
		top.ColorScale.ScaleWithColor(colorHighlight)
		text.Draw(screen, fmt.Sprintf("Рекорд: %d", g.highScore), g.getFontFace(18), top)
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

// drawCell рисует клетку на сетке
func drawCell(screen *ebiten.Image, pos entity.Point, c color.RGBA) {
	x := float32(pos.X * CellSize)
	y := float32(pos.Y*CellSize + CellSize + 4)

	// Оптимизировано: рисуем основное тело клетки одним вызовом
	vector.FillRect(screen, x, y, float32(CellSize), float32(CellSize), c, false)

	// Рисуем границы клетки
	borderColor := color.RGBA{255, 255, 255, 30}
	// Верхняя граница
	vector.FillRect(screen, x, y, float32(CellSize), 1, borderColor, false)
	// Нижняя граница
	vector.FillRect(screen, x, y+float32(CellSize)-1, float32(CellSize), 1, borderColor, false)
	// Левая граница
	vector.FillRect(screen, x, y, 1, float32(CellSize), borderColor, false)
	// Правая граница
	vector.FillRect(screen, x+float32(CellSize)-1, y, 1, float32(CellSize), borderColor, false)
}
