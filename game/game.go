package game

import (
	"image/color"
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

// Layout возвращает размер игрового экрана
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	b := g.board
	return b.ScreenWidth(), b.Y + b.Height*CellSize + 4
}