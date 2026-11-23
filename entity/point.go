package entity

// GridSize определяет размер игрового поля в клетках
const (
	GridWidth  = 30
	GridHeight = 25

	// Начальная длина змейки
	StartSnakeLength = 3
)

// Point представляет точку на сетке
type Point struct {
	X, Y int
}

// Equals проверяет равенство двух точек
func (p Point) Equals(other Point) bool {
	return p.X == other.X && p.Y == other.Y
}

// Add возвращает новую точку, сдвинутую на dx, dy
func (p Point) Add(dx, dy int) Point {
	return Point{X: p.X + dx, Y: p.Y + dy}
}
