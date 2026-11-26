package game

import "github.com/rekhin/snake-game/entity"

// HUDHeight — высота панели интерфейса над игровым полем в пикселях
const HUDHeight = 40

// Board представляет игровое поле (сетку) с координатами и размерами.
// Все объекты змейки, еды и эффектов позиционируются относительно поля.
type Board struct {
	X      int // позиция верхнего левого угла сетки по X (экранная координата)
	Y      int // позиция верхнего левого угла сетки по Y (экранная координата)
	Width  int // количество клеток по горизонтали
	Height int // количество клеток по вертикали
}

// NewBoard создаёт игровое поле стандартного размера.
func NewBoard() Board {
	return Board{
		X:      0,
		Y:      HUDHeight,
		Width:  entity.GridWidth,
		Height: entity.GridHeight,
	}
}

// PixelX возвращает экранный X-пиксель для заданной колонки сетки.
func (b Board) PixelX(col int) int {
	return b.X + col*CellSize
}

// PixelY возвращает экранный Y-пиксель для заданной строки сетки.
func (b Board) PixelY(row int) int {
	return b.Y + row*CellSize
}

// GridToPixel преобразует координаты клетки в экранные пиксели (верхний левый угол).
func (b Board) GridToPixel(col, row int) (float32, float32) {
	return float32(b.PixelX(col)), float32(b.PixelY(row))
}

// GridToPixelCenter преобразует координаты клетки в экранные пиксели центра.
func (b Board) GridToPixelCenter(col, row int) (float32, float32) {
	return float32(b.PixelX(col) + CellSize/2), float32(b.PixelY(row) + CellSize/2)
}

// Bounds возвращает прямоугольник поля в экранных координатах.
func (b Board) Bounds() (int, int, int, int) {
	return b.X, b.Y, b.Width * CellSize, b.Height * CellSize
}

// ScreenHeight возвращает полную высоту экрана с учётом HUD и игрового поля.
func (b Board) ScreenHeight() int {
	return b.Y + b.Height*CellSize + 4
}

// ScreenWidth возвращает ширину экрана (совпадает с шириной поля).
func (b Board) ScreenWidth() int {
	return b.Width * CellSize
}