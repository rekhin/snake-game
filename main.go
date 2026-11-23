package main

import (
	_ "embed"
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/rekhin/snake-game/game"
)

//go:embed resources/NotoSans-Regular.ttf
var fontData []byte

const (
	windowWidth  = game.GridWidth * game.CellSize
	windowHeight = game.GridHeight*game.CellSize + game.CellSize + 4
	windowTitle  = "🐍 Змейка"
)

func main() {
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle(windowTitle)
	ebiten.SetWindowResizable(true)

	g := game.NewGame(fontData)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Println("Игра завершена")
}
