package entity

import (
	"math/rand"
	"time"
)

// Food представляет еду на поле
type Food struct {
	Position Point
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// NewFood создаёт новую еду в случайной позиции
func NewFood(snakeBody []Point) *Food {
	f := &Food{}
	f.Spawn(snakeBody)
	return f
}

// Spawn размещает еду в случайной пустой клетке
func (f *Food) Spawn(snakeBody []Point) {
	occupied := make(map[Point]bool)
	for _, p := range snakeBody {
		occupied[p] = true
	}

	for {
		x := rand.Intn(GridWidth)
		y := rand.Intn(GridHeight)
		pos := Point{X: x, Y: y}
		if !occupied[pos] {
			f.Position = pos
			return
		}
	}
}
