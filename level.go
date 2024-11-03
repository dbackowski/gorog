package main

import (
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Level struct {
	Tiles []MapTile
}

func NewLevel() Level {
	l := Level{}
	tiles := l.CreateTiles()
	l.Tiles = tiles
	return l
}

func (level *Level) GetIndexFromXY(x int, y int) int {
	gd := NewGameData()
	return (y * gd.ScreenWidth) + x
}

func (level *Level) CreateTiles() []MapTile {
	gd := NewGameData()
	tiles := make([]MapTile, 0)
	img, _, err := ebitenutil.NewImageFromFile("assets/EverRogueTileset 1.0 Horizontal.png")

	if err != nil {
		log.Fatal(err)
	}

	for x := 0; x < gd.ScreenWidth; x++ {
		for y := 0; y < gd.ScreenHeight; y++ {
			if x == 0 || x == gd.ScreenWidth-1 || y == 0 || y == gd.ScreenHeight-1 {
				wall := img.SubImage(image.Rect(272, 0, 256, 272)).(*ebiten.Image)

				tile := MapTile{
					PixelX:  x * gd.TileWidth,
					PixelY:  y * gd.TileHeight,
					Blocked: true,
					Image:   wall,
				}
				tiles = append(tiles, tile)
			} else {
				floor := img.SubImage(image.Rect(0, 0, 16, 16)).(*ebiten.Image)

				tile := MapTile{
					PixelX:  x * gd.TileWidth,
					PixelY:  y * gd.TileHeight,
					Blocked: false,
					Image:   floor,
				}
				tiles = append(tiles, tile)
			}
		}
	}

	return tiles
}
