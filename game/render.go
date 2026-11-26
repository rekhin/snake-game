package game

import (
	"fmt"
	"image/color"
	"math"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

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

	titleText := "🐍 ЗМЕЙКА"
	titleFace := g.getFontFace(48)
	titleW, _ := text.Measure(titleText, titleFace, 0)
	top := &text.DrawOptions{}
	top.GeoM.Translate(float64(w/2)-titleW/2, float64(h/3))
	top.ColorScale.ScaleWithColor(colorTitle)
	text.Draw(screen, titleText, titleFace, top)

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

	op := &ebiten.DrawImageOptions{}
	if !g.shakeEndTime.IsZero() {
		op.GeoM.Translate(float64(g.shakeOffset[0]), float64(g.shakeOffset[1]))
	}
	screen.DrawImage(g.gridCache, op)

	head := g.snake.GetHead()
	var borderColor color.RGBA
	for i := 1; i < len(g.snake.Body); i++ {
		pos := g.snake.Body[i]
		dx := float32(pos.X * CellSize)
		dy := float32(pos.Y*CellSize + CellSize + 4)
		vector.FillRect(screen, dx, dy, float32(CellSize), float32(CellSize), colorSnakeBody, false)

		borderColor = color.RGBA{255, 255, 255, 20}
		vector.StrokeLine(screen, dx, dy, dx+float32(CellSize), dy, 1, borderColor, false)
		vector.StrokeLine(screen, dx, dy+float32(CellSize)-1, dx+float32(CellSize), dy+float32(CellSize)-1, 1, borderColor, false)
		vector.StrokeLine(screen, dx, dy, dx, dy+float32(CellSize), 1, borderColor, false)
		vector.StrokeLine(screen, dx+float32(CellSize)-1, dy, dx+float32(CellSize)-1, dy+float32(CellSize), 1, borderColor, false)
	}

	headX := float32(head.X * CellSize)
	headY := float32(head.Y*CellSize + CellSize + 4)
	vector.FillRect(screen, headX, headY, float32(CellSize), float32(CellSize), colorSnakeHead, false)

	borderColor = color.RGBA{255, 255, 255, 40}
	vector.StrokeLine(screen, headX, headY, headX+float32(CellSize), headY, 1.5, borderColor, false)
	vector.StrokeLine(screen, headX, headY+float32(CellSize)-1, headX+float32(CellSize), headY+float32(CellSize)-1, 1.5, borderColor, false)
	vector.StrokeLine(screen, headX, headY, headX, headY+float32(CellSize), 1.5, borderColor, false)
	vector.StrokeLine(screen, headX+float32(CellSize)-1, headY, headX+float32(CellSize)-1, headY+float32(CellSize), 1.5, borderColor, false)

	foodPos := g.food.Position
	foodCenterX := float32(foodPos.X*CellSize + CellSize/2)
	foodCenterY := float32(foodPos.Y*CellSize+CellSize/2 + CellSize + 4)

	var foodColor = colorFood
	if g.pulsePhase > 0 {
		scale := 1.0 + math.Sin(g.pulsePhase)*0.15
		radius := float32(CellSize/2-2) * float32(scale)

		glowColor := color.RGBA{235, 87, 87, 80}
		vector.DrawFilledCircle(screen, foodCenterX, foodCenterY, radius+4, glowColor, false)

		foodColor = color.RGBA{
			uint8(235 * (0.8 + math.Sin(g.pulsePhase)*0.2)),
			uint8(87 * (0.8 + math.Sin(g.pulsePhase)*0.2)),
			uint8(87 * (0.8 + math.Sin(g.pulsePhase)*0.2)),
			255,
		}
		vector.DrawFilledCircle(screen, foodCenterX, foodCenterY, radius, foodColor, false)

		highlight := color.RGBA{255, 180, 180, 120}
		vector.DrawFilledCircle(screen, foodCenterX-2, foodCenterY-2, radius*0.3, highlight, false)
	} else {
		radius := float32(CellSize/2 - 2)
		vector.DrawFilledCircle(screen, foodCenterX, foodCenterY, radius, colorFood, false)

		highlight := color.RGBA{255, 180, 180, 120}
		vector.DrawFilledCircle(screen, foodCenterX-2, foodCenterY-2, radius*0.3, highlight, false)
	}

	g.drawHUD(screen)
}

// drawOverlay рисует полупрозрачный оверлей с текстом
func (g *Game) drawOverlay(screen *ebiten.Image, overlayText string) {
	w, h := screen.Bounds().Dx(), screen.Bounds().Dy()

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