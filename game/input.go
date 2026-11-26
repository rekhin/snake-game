package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/rekhin/snake-game/entity"
)

// processInputPause обрабатывает клавиши паузы и возвращает true, если была нажата клавиша паузы
func (g *Game) processInputPause() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.state = StatePaused
		return true
	}
	return false
}

// processInputDirection обрабатывает клавиши направления змейки
func (g *Game) processInputDirection() {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		g.inputQueue = append(g.inputQueue, entity.DirUp)
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.inputQueue = append(g.inputQueue, entity.DirDown)
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
		g.inputQueue = append(g.inputQueue, entity.DirLeft)
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.inputQueue = append(g.inputQueue, entity.DirRight)
	}

	// Ограничение очереди (максимум 2 ввода для отзывчивости без задержек)
	if len(g.inputQueue) > 2 {
		g.inputQueue = g.inputQueue[len(g.inputQueue)-2:]
	}

	// Применение одного направления из очереди за один тик
	if len(g.inputQueue) > 0 {
		nextDir := g.inputQueue[0]
		g.inputQueue = g.inputQueue[1:]
		g.snake.SetDirection(nextDir)
	}
}

// processInput обрабатывает клавиши старта игры и возвращает true, если игра началась
func (g *Game) processInput() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		return true
	}
	return false
}