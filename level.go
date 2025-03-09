package main

import (
	"fmt"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/norendren/go-fov/fov"
)

var floor *ebiten.Image
var wall *ebiten.Image
var levelHeight int = 0

type TileType int

const (
	WALL TileType = iota
	FLOOR
	STAIRS_DOWN
	STAIRS_UP
)

type Level struct {
	Tiles         []*MapTile
	Rooms         []Rect
	PlayerVisible *fov.View
	StairsDown    *Position
	StairsUp      *Position
	Depth         int
}

var stairsDown *ebiten.Image
var stairsUp *ebiten.Image

func loadTileImages() {
	if floor != nil && wall != nil && stairsDown != nil && stairsUp != nil {
		return
	}

	var err error

	floor, _, err = ebitenutil.NewImageFromFile("assets/floor.png")
	if err != nil {
		log.Fatal(err)
	}

	wall, _, err = ebitenutil.NewImageFromFile("assets/wall.png")
	if err != nil {
		log.Fatal(err)
	}

	// For now, we'll reuse the floor image for stairs but with different colors
	stairsDown, _, err = ebitenutil.NewImageFromFile("assets/stairs_down.png")
	if err != nil {
		log.Fatal(err)
	}

	stairsUp, _, err = ebitenutil.NewImageFromFile("assets/stairs_up.png")
	if err != nil {
		log.Fatal(err)
	}
}

// Max returns the larger of x or y.
func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

// Min returns the smaller of x or y.
func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func (level Level) InBounds(x, y int) bool {
	gd := NewGameData()
	if x < 0 || x > gd.ScreenWidth || y < 0 || y > levelHeight {
		return false
	}
	return true
}

func (level Level) IsOpaque(x, y int) bool {
	idx := level.GetIndexFromXY(x, y)
	return level.Tiles[idx].TileType == WALL
}

func (level *Level) createHorizontalTunnel(x1 int, x2 int, y int) {
	gd := NewGameData()

	for x := min(x1, x2); x < max(x1, x2)+1; x++ {
		index := level.GetIndexFromXY(x, y)

		if index > 0 && index < gd.ScreenWidth*levelHeight {
			level.Tiles[index].Blocked = false
			level.Tiles[index].TileType = FLOOR
			level.Tiles[index].Image = floor
		}
	}
}

func (level *Level) createVerticalTunnel(y1 int, y2 int, x int) {
	gd := NewGameData()

	for y := min(y1, y2); y < max(y1, y2)+1; y++ {
		index := level.GetIndexFromXY(x, y)

		if index > 0 && index < gd.ScreenWidth*levelHeight {
			level.Tiles[index].Blocked = false
			level.Tiles[index].TileType = FLOOR
			level.Tiles[index].Image = floor
		}
	}
}

func (level *Level) createRoom(room Rect) {
	for y := room.Y1 + 1; y < room.Y2; y++ {
		for x := room.X1 + 1; x < room.X2; x++ {
			index := level.GetIndexFromXY(x, y)
			level.Tiles[index].Blocked = false
			level.Tiles[index].TileType = FLOOR
			level.Tiles[index].Image = floor
		}
	}
}

func (level *Level) GenerateLevelTiles() {
	MIN_SIZE := 6
	MAX_SIZE := 10
	MAX_ROOMS := 30
	contains_rooms := false

	gd := NewGameData()
	levelHeight = gd.ScreenHeight - gd.UIHeight
	tiles := level.createTiles()
	level.Tiles = tiles

	for idx := 0; idx < MAX_ROOMS; idx++ {
		w := GetRandomBetween(MIN_SIZE, MAX_SIZE)
		h := GetRandomBetween(MIN_SIZE, MAX_SIZE)
		x := GetDiceRoll(gd.ScreenWidth - w - 1)
		y := GetDiceRoll(levelHeight - h - 1)

		new_room := NewRect(x, y, w, h)
		okToAdd := true
		for _, otherRoom := range level.Rooms {
			if new_room.Intersect(otherRoom) {
				okToAdd = false
				break
			}
		}
		if okToAdd {
			level.createRoom(new_room)
			if contains_rooms {
				newX, newY := new_room.Center()
				prevX, prevY := level.Rooms[len(level.Rooms)-1].Center()
				coinflip := GetDiceRoll(2)
				if coinflip == 2 {
					level.createHorizontalTunnel(prevX, newX, prevY)
					level.createVerticalTunnel(prevY, newY, newX)

				} else {
					level.createHorizontalTunnel(prevX, newX, newY)
					level.createVerticalTunnel(prevY, newY, prevX)
				}

			}

			level.Rooms = append(level.Rooms, new_room)
			contains_rooms = true
		}
	}

	// Place stairs down in a random room (not the first room where player starts)
	if len(level.Rooms) > 1 {
		// Choose a random room (not the first one)
		roomIdx := GetDiceRoll(len(level.Rooms) - 1)
		if roomIdx == 0 {
			roomIdx = 1 // Ensure we don't use the first room
		}

		// Get center of the room
		x, y := level.Rooms[roomIdx].Center()

		// Place stairs down
		level.StairsDown = &Position{X: x, Y: y}
		fmt.Printf("Stairs down at %d, %d\n", x, y)
		idx := level.GetIndexFromXY(x, y)
		level.Tiles[idx].TileType = STAIRS_DOWN
		level.Tiles[idx].Image = stairsDown
	}

	// If this is not the first level, place stairs up
	if level.Depth > 0 {
		// Place stairs up in the first room
		x, y := level.Rooms[0].Center()

		// Place stairs up
		level.StairsUp = &Position{X: x, Y: y}
		idx := level.GetIndexFromXY(x, y)
		level.Tiles[idx].TileType = STAIRS_UP
		level.Tiles[idx].Image = stairsUp
	}
}

func NewLevel(depth int) Level {
	l := Level{
		Depth: depth,
	}
	loadTileImages()
	rooms := make([]Rect, 0)
	l.Rooms = rooms
	l.GenerateLevelTiles()
	l.PlayerVisible = fov.New()
	return l
}

func (level *Level) DrawLevel(screen *ebiten.Image, debugMode bool) {
	gd := NewGameData()

	for x := 0; x < gd.ScreenWidth; x++ {
		for y := 0; y < levelHeight; y++ {
			idx := level.GetIndexFromXY(x, y)
			tile := level.Tiles[idx]
			isVis := debugMode || level.PlayerVisible.IsVisible(x, y)

			if isVis {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(tile.PixelX), float64(tile.PixelY))
				screen.DrawImage(tile.Image, op)
				level.Tiles[idx].IsRevealed = true
			} else if tile.IsRevealed {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(tile.PixelX), float64(tile.PixelY))
				op.ColorM.Translate(-0.3, -0.3, -0.3, 0.35)
				screen.DrawImage(tile.Image, op)
			}
		}
	}
}

func (level *Level) GetIndexFromXY(x int, y int) int {
	gd := NewGameData()
	return (y * gd.ScreenWidth) + x
}

func (level *Level) createTiles() []*MapTile {
	gd := NewGameData()
	tiles := make([]*MapTile, levelHeight*gd.ScreenWidth)
	index := 0
	for x := 0; x < gd.ScreenWidth; x++ {
		for y := 0; y < levelHeight; y++ {
			index = level.GetIndexFromXY(x, y)
			tile := MapTile{
				PixelX:     x * gd.TileWidth,
				PixelY:     y * gd.TileHeight,
				Blocked:    true,
				Image:      wall,
				IsRevealed: false,
				TileType:   WALL,
			}
			tiles[index] = &tile
		}
	}
	return tiles
}
