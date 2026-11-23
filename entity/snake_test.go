package entity

import (
	"testing"
)

func TestNewSnake_DefaultPosition(t *testing.T) {
	snake := NewSnake()

	if len(snake.Body) != StartSnakeLength {
		t.Errorf("expected snake length %d, got %d", StartSnakeLength, len(snake.Body))
	}

	head := snake.GetHead()
	if head.X != GridWidth/2 || head.Y != GridHeight/2 {
		t.Errorf("expected head at (%d,%d), got (%d,%d)", GridWidth/2, GridHeight/2, head.X, head.Y)
	}
}

func TestNewSnake_DefaultDirection(t *testing.T) {
	snake := NewSnake()

	if snake.Direction != DirRight {
		t.Errorf("expected direction Right, got %v", snake.Direction)
	}
}

func TestSetDirection_ChangeDirection(t *testing.T) {
	snake := NewSnake()

	snake.SetDirection(DirUp)
	if snake.Direction != DirUp {
		t.Error("expected direction to be Up")
	}

	snake.SetDirection(DirLeft)
	if snake.Direction != DirLeft {
		t.Error("expected direction to be Left")
	}
}

func TestSetDirection_OppositeBlocked(t *testing.T) {
	snake := NewSnake() // starts moving Right

	// Trying to set opposite direction (Left) should be blocked
	result := snake.SetDirection(DirLeft)
	if result {
		t.Error("opposite direction change should return false")
	}
	if snake.Direction != DirRight {
		t.Errorf("expected direction to stay Right, got %v", snake.Direction)
	}

	snake2 := NewSnake()
	snake2.SetDirection(DirUp)

	// Trying to set opposite (Down) should be blocked
	result = snake2.SetDirection(DirDown)
	if result {
		t.Error("opposite direction change should return false")
	}
	if snake2.Direction != DirUp {
		t.Errorf("expected direction to stay Up, got %v", snake2.Direction)
	}
}

func TestMove_MoveRight(t *testing.T) {
	snake := NewSnake()
	startX := snake.GetHead().X

	snake.SetDirection(DirRight)
	snake.Move()

	head := snake.GetHead()
	if head.X != startX+1 {
		t.Errorf("expected x to increase by 1, got %d (was %d)", head.X, startX)
	}
	if head.Y != GridHeight/2 { // Y should stay the same
		t.Errorf("expected y to stay %d, got %d", GridHeight/2, head.Y)
	}
}

func TestMove_MoveUp(t *testing.T) {
	snake := NewSnake()
	startY := snake.GetHead().Y

	snake.SetDirection(DirUp)
	snake.Move()

	head := snake.GetHead()
	if head.Y != startY-1 {
		t.Errorf("expected y to decrease by 1, got %d (was %d)", head.Y, startY)
	}
}

func TestMove_MoveDown(t *testing.T) {
	snake := NewSnake()
	startY := snake.GetHead().Y

	snake.SetDirection(DirDown)
	snake.Move()

	head := snake.GetHead()
	if head.Y != startY+1 {
		t.Errorf("expected y to increase by 1, got %d (was %d)", head.Y, startY)
	}
}

func TestMove_MoveLeft(t *testing.T) {
	snake := NewSnake()
	startX := snake.GetHead().X

	// First move up to change direction, then left
	snake.SetDirection(DirUp)
	snake.Move()
	startX = snake.GetHead().X // new head after moving up

	snake.SetDirection(DirLeft)
	snake.Move()

	head := snake.GetHead()
	if head.X != startX-1 {
		t.Errorf("expected x to decrease by 1, got %d (was %d)", head.X, startX)
	}
}

func TestGrowBy_IncreasesLength(t *testing.T) {
	snake := NewSnake()
	initialLen := len(snake.Body)

	// GrowBy sets a counter; actual growth happens on next Move()
	snake.GrowBy(3)

	// After 1 move with grow=3, the snake should have grown by 1 (no tail removed)
	snake.Move()
	if len(snake.Body) != initialLen+1 {
		t.Errorf("after 1st move: expected length %d, got %d", initialLen+1, len(snake.Body))
	}

	snake.Move()
	if len(snake.Body) != initialLen+2 {
		t.Errorf("after 2nd move: expected length %d, got %d", initialLen+2, len(snake.Body))
	}

	snake.Move()
	if len(snake.Body) != initialLen+3 {
		t.Errorf("after 3rd move: expected length %d, got %d", initialLen+3, len(snake.Body))
	}

	// After grow is exhausted, normal movement should not increase length
	snake.Move()
	if len(snake.Body) != initialLen+3 {
		t.Errorf("after 4th move (grow exhausted): expected length %d, got %d", initialLen+3, len(snake.Body))
	}
}

func TestGrowBy_Zero(t *testing.T) {
	snake := NewSnake()
	initialLen := len(snake.Body)

	snake.GrowBy(0)
	snake.Move() // one move after grow=0 should not change length
	if len(snake.Body) != initialLen {
		t.Errorf("expected length to stay %d, got %d", initialLen, len(snake.Body))
	}
}

func TestCheckWallCollision_OutOfBounds(t *testing.T) {
	snake := NewSnake()

	// Move snake to the edge and then beyond (right boundary)
	for i := 0; i < GridWidth+5; i++ {
		snake.SetDirection(DirRight)
		snake.Move()
	}

	if !snake.CheckWallCollision() {
		t.Error("expected wall collision when x >= GridWidth")
	}

	// Test left boundary
	snake2 := NewSnake()
	for i := 0; i < GridWidth/2+5; i++ {
		snake2.SetDirection(DirLeft)
		snake2.Move()
	}

	if !snake2.CheckWallCollision() {
		t.Error("expected wall collision when x < 0")
	}

	// Test top boundary
	snake3 := NewSnake()
	for i := 0; i < GridHeight/2+5; i++ {
		snake3.SetDirection(DirUp)
		snake3.Move()
	}

	if !snake3.CheckWallCollision() {
		t.Error("expected wall collision when y < 0")
	}

	// Test bottom boundary
	snake4 := NewSnake()
	for i := 0; i < GridHeight/2+5; i++ {
		snake4.SetDirection(DirDown)
		snake4.Move()
	}

	if !snake4.CheckWallCollision() {
		t.Error("expected wall collision when y >= GridHeight")
	}
}

func TestCheckWallCollision_WithinBounds(t *testing.T) {
	snake := NewSnake()

	if snake.CheckWallCollision() {
		t.Error("new snake should not collide with walls")
	}

	// Move within bounds
	for i := 0; i < 5; i++ {
		snake.SetDirection(DirRight)
		snake.Move()
	}

	if snake.CheckWallCollision() {
		t.Error("snake should not collide with walls when still in bounds")
	}
}

func TestCheckSelfCollision_NoCollision(t *testing.T) {
	snake := NewSnake()

	if snake.CheckSelfCollision() {
		t.Error("new snake should not self-collide")
	}
}

func TestIsBody_ReturnsCorrectly(t *testing.T) {
	snake := NewSnake()
	head := snake.GetHead()

	if snake.IsBody(head) {
		t.Error("head should not be considered body")
	}

	snake.GrowBy(1)
	bodyPart := snake.Body[1]

	if !snake.IsBody(bodyPart) {
		t.Errorf("expected (%d,%d) to be part of body", bodyPart.X, bodyPart.Y)
	}
}

func TestPoint_Equals(t *testing.T) {
	p1 := Point{X: 5, Y: 10}
	p2 := Point{X: 5, Y: 10}
	p3 := Point{X: 6, Y: 10}

	if !p1.Equals(p2) {
		t.Error("same points should be equal")
	}

	if p1.Equals(p3) {
		t.Error("different points should not be equal")
	}
}

func TestPoint_Add(t *testing.T) {
	p := Point{X: 5, Y: 10}
	result := p.Add(2, -3)

	if result.X != 7 || result.Y != 7 {
		t.Errorf("expected (7,7), got (%d,%d)", result.X, result.Y)
	}
}

func TestNewFood_SpawnsInEmptyCell(t *testing.T) {
	snake := NewSnake()
	food := NewFood(snake.Body)

	if food.Position.Equals(snake.GetHead()) {
		t.Error("food should not spawn on snake head")
	}
}

func TestNewFood_SpawnsWithinBounds(t *testing.T) {
	snake := NewSnake()
	food := NewFood(snake.Body)

	if food.Position.X < 0 || food.Position.X >= GridWidth {
		t.Errorf("food x=%d out of bounds [0,%d)", food.Position.X, GridWidth)
	}
	if food.Position.Y < 0 || food.Position.Y >= GridHeight {
		t.Errorf("food y=%d out of bounds [0,%d)", food.Position.Y, GridHeight)
	}
}

func TestFood_SpawnAvoidsSnake(t *testing.T) {
	snake := NewSnake()
	// Grow snake to fill most of the board
	for i := 0; i < 100; i++ {
		snake.GrowBy(1)
	}

	food := NewFood(snake.Body)

	for _, segment := range snake.Body {
		if food.Position.Equals(segment) {
			t.Errorf("food spawned on snake body at (%d,%d)", segment.X, segment.Y)
		}
	}
}