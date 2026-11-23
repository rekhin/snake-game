package game

import "github.com/rekhin/snake-game/entity"

// Константы конфигурации игры (используются из entity)
const (
	// Размер клетки в пикселях
	CellSize = 20

	// Начальная скорость (мс между ходами)
	StartMoveIntervalMs = 150

	// Минимальная скорость (мс между ходами)
	MinMoveIntervalMs = 60

	// Уменьшение интервала каждые N очков
	SpeedUpEveryNPoints = 5

	// Награда за еду в очках
	FoodScore = 10
)

// Алиасы для удобства использования
const (
	GridWidth  = entity.GridWidth
	GridHeight = entity.GridHeight
)
