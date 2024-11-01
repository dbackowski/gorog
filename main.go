package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct{}

func NewGame() *Game {
	g := &Game{}
	return g
}

// Update is called each tic.
func (g *Game) Update() error {
	return nil
}

// Draw is called each draw cycle and is where we will blit.
func (g *Game) Draw(screen *ebiten.Image) {
}

// Layout will return the screen dimensions.
func (g *Game) Layout(w, h int) (int, int) { return 1280, 800 }

func main() {
	g := NewGame()
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)
	ebiten.SetWindowTitle("Gorog")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
