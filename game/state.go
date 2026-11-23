package game

// State представляет состояние игры
type State int

const (
	StateMenu State = iota // Главное меню
	StatePlaying           // Игра идёт
	StatePaused            // Пауза
	StateGameOver          // Конец игры
)

func (s State) String() string {
	switch s {
	case StateMenu:
		return "MENU"
	case StatePlaying:
		return "PLAYING"
	case StatePaused:
		return "PAUSED"
	case StateGameOver:
		return "GAME_OVER"
	default:
		return "UNKNOWN"
	}
}
