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

func main() {
	g := game.NewGame(fontData)

	windowW, windowH := g.Layout(0, 0)
	ebiten.SetWindowSize(windowW, windowH)
	ebiten.SetWindowTitle("🐍 Змейка")
	ebiten.SetWindowResizable(true)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Println("Игра завершена")
}
