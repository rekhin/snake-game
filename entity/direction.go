package entity

// Direction представляет направление движения змейки
type Direction int

const (
	DirUp Direction = iota
	DirRight
	DirDown
	DirLeft
	DirNone // Начальное состояние, нет направления
)

// Opposite возвращает противоположное направление
func (d Direction) Opposite() Direction {
	switch d {
	case DirUp:
		return DirDown
	case DirDown:
		return DirUp
	case DirLeft:
		return DirRight
	case DirRight:
		return DirLeft
	default:
		return DirNone
	}
}

// IsOpposite проверяет, являются ли направления противоположными
func (d Direction) IsOpposite(other Direction) bool {
	return d == other.Opposite()
}
