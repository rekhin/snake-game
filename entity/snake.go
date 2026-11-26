package entity

// Snake представляет змейку
type Snake struct {
	Body      []Point // Тела змейки (голова — Body[0])
	Direction Direction // Текущее направление
	Grow      int       // Сколько сегментов ещё нужно добавить
}

// NewSnake создаёт новую змейку в центре поля
func NewSnake() *Snake {
	body := make([]Point, StartSnakeLength)
	centerX, centerY := GridWidth/2, GridHeight/2
	for i := 0; i < StartSnakeLength; i++ {
		body[i] = Point{X: centerX - i, Y: centerY}
	}
	return &Snake{
		Body:      body,
		Direction: DirRight,
		Grow:      0,
	}
}

// Move двигает змейку на один шаг в текущем направлении.
// Если Grow > 0, хвост не удаляется (змейка растёт).
func (s *Snake) Move() {
	if len(s.Body) == 0 {
		return
	}

	head := s.Body[0]
	var newHead Point
	switch s.Direction {
	case DirUp:
		newHead = head.Add(0, -1)
	case DirDown:
		newHead = head.Add(0, 1)
	case DirLeft:
		newHead = head.Add(-1, 0)
	case DirRight:
		newHead = head.Add(1, 0)
	default:
		return // Нет направления — не двигаемся
	}

	// Добавляем новую голову в начало
	s.Body = append([]Point{newHead}, s.Body...)

	// Если нужно расти, не удаляем хвост
	if s.Grow > 0 {
		s.Grow--
		return
	}

	// Удаляем хвост
	s.Body = s.Body[:len(s.Body)-1]
}

// GrowBy увеличивает змейку на n сегментов
func (s *Snake) GrowBy(n int) {
	s.Grow += n
}

// SetDirection устанавливает новое направление с проверкой.
// Возвращает false, если направление противоположно текущему или равно DirNone.
func (s *Snake) SetDirection(newDir Direction) bool {
	if newDir == DirNone || s.Direction.IsOpposite(newDir) {
		return false // Нельзя двигаться в противоположную сторону
	}
	s.Direction = newDir
	return true
}

// CheckWallCollision проверяет столкновение со стенами
func (s *Snake) CheckWallCollision() bool {
	head := s.Body[0]
	return head.X < 0 || head.X >= GridWidth || head.Y < 0 || head.Y >= GridHeight
}

// CheckSelfCollision проверяет столкновение с собой
func (s *Snake) CheckSelfCollision() bool {
	head := s.Body[0]
	for i := 1; i < len(s.Body); i++ {
		if head.Equals(s.Body[i]) {
			return true
		}
	}
	return false
}

// GetHead возвращает позицию головы
func (s *Snake) GetHead() Point {
	if len(s.Body) == 0 {
		return Point{}
	}
	return s.Body[0]
}

// IsBody проверяет, находится ли точка в теле змейки (не считая голову)
func (s *Snake) IsBody(p Point) bool {
	for i := 1; i < len(s.Body); i++ {
		if s.Body[i].Equals(p) {
			return true
		}
	}
	return false
}
