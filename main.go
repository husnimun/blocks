package main

import (
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	board Board
}

func (g *Game) Update() error {
	g.board.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.board.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 167, 210

}

func main() {
	ebiten.SetWindowSize(334, 420)
	ebiten.SetWindowTitle("Blocks")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
