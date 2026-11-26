package entity

import (
	"math/rand"
	"testing"
)

// TestFood_SpawnDeterministic проверяет, что с фиксированным seed спавн еды
// всегда возвращает одинаковую позицию — это гарантирует детерминизм тестов.
func TestFood_SpawnDeterministic(t *testing.T) {
	// Сохраняем оригинальный источник и восстанавливаем его после теста
	originalSource := randSource

	snakeBody := []Point{{X: 5, Y: 5}, {X: 6, Y: 5}, {X: 7, Y: 5}}

	// Создаём еду с seed=12345
	SetRandSource(rand.New(rand.NewSource(12345)))
	food1 := NewFood(snakeBody)
	pos1 := food1.Position

	// Восстанавливаем тот же seed и создаём вторую еду
	SetRandSource(rand.New(rand.NewSource(12345)))
	food2 := NewFood(snakeBody)
	pos2 := food2.Position

	// Восстанавливаем оригинальный источник
	SetRandSource(originalSource)

	if pos1 != pos2 {
		t.Errorf("non-deterministic spawn: first=%v, second=%v (expected equal with fixed seed)",
			pos1, pos2)
	}

	// Проверяем, что позиция в допустимых границах и не пересекается со змейкой
	if pos1.X < 0 || pos1.X >= GridWidth {
		t.Errorf("spawn X=%d out of bounds [0, %d)", pos1.X, GridWidth)
	}
	if pos1.Y < 0 || pos1.Y >= GridHeight {
		t.Errorf("spawn Y=%d out of bounds [0, %d)", pos1.Y, GridHeight)
	}

	occupied := make(map[Point]bool)
	for _, p := range snakeBody {
		occupied[p] = true
	}
	if occupied[pos1] {
		t.Errorf("spawn position %v overlaps with snake body", pos1)
	}
}

// TestFood_SpawnWithLargeSnake_Deterministic проверяет детерминизм при оптимизированном
// спавне (когда змейка занимает больше половины поля).
func TestFood_SpawnWithLargeSnake_Deterministic(t *testing.T) {
	originalSource := randSource
	snakeBody := make([]Point, 400)
	for i := range snakeBody {
		snakeBody[i] = Point{X: i % GridWidth, Y: (i / GridWidth) % GridHeight}
	}

	// Создаём еду с seed=99999
	SetRandSource(rand.New(rand.NewSource(99999)))
	food1 := NewFood(snakeBody)
	pos1 := food1.Position

	// Восстанавливаем тот же seed и создаём вторую еду
	SetRandSource(rand.New(rand.NewSource(99999)))
	food2 := NewFood(snakeBody)
	pos2 := food2.Position

	// Восстанавливаем оригинальный источник
	SetRandSource(originalSource)

	if pos1 != pos2 {
		t.Errorf("non-deterministic spawn (large snake): first=%v, second=%v",
			pos1, pos2)
	}

	// Проверяем отсутствие коллизии со змейкой
	occupied := make(map[Point]bool)
	for _, p := range snakeBody {
		occupied[p] = true
	}
	if occupied[pos1] {
		t.Errorf("spawn position %v overlaps with large snake body", pos1)
	}
}

// TestFood_SpawnRandomness проверяет, что разные seed дают разные результаты.
func TestFood_SpawnRandomness(t *testing.T) {
	originalSource := randSource
	snakeBody := []Point{{X: 5, Y: 5}, {X: 6, Y: 5}}

	results := make(map[Point]int)
	for i := 0; i < 100; i++ {
		SetRandSource(rand.New(rand.NewSource(int64(i))))
		food := NewFood(snakeBody)
		results[food.Position]++
	}

	// Восстанавливаем оригинальный источник
	SetRandSource(originalSource)

	// С 100 разными seeds должны получить несколько разных позиций
	if len(results) < 5 {
		t.Errorf("expected at least 5 different positions from 100 different seeds, got %d", len(results))
	}
}
