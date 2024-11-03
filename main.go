package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	Map GameMap
}

type MapTile struct {
	PixelX  int
	PixelY  int
	Blocked bool
	Image   *ebiten.Image
}

func NewGame() *Game {
	g := &Game{}
	g.Map = NewGameMap()
	return g
}

func (g *Game) Layout(w, h int) (int, int) {
	gd := NewGameData()
	return gd.TileWidth * gd.ScreenWidth, gd.TileHeight * gd.ScreenHeight
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	gd := NewGameData()
	level := g.Map.Dungeons[0].Levels[0]

	for x := 0; x < gd.ScreenWidth; x++ {
		for y := 0; y < gd.ScreenHeight; y++ {
			tile := level.Tiles[level.GetIndexFromXY(x, y)]
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(tile.PixelX), float64(tile.PixelY))
			screen.DrawImage(tile.Image, op)
		}
	}
}

func main() {
	g := NewGame()
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)
	ebiten.SetWindowTitle("Gorog")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
