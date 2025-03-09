package main

import (
	"flag"
	"log"

	"github.com/bytearena/ecs"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	Map         GameMap
	World       *ecs.Manager
	WorldTags   map[string]ecs.Tag
	Turn        TurnState
	TurnCounter int
	DebugMode   bool
}
type MapTile struct {
	PixelX     int
	PixelY     int
	Blocked    bool
	Image      *ebiten.Image
	IsRevealed bool
	TileType   TileType
}

func NewGame(debugMode bool) *Game {
	g := &Game{}
	g.Map = NewGameMap()
	world, tags := InitializeWorld(g.Map.CurrentLevel)
	g.WorldTags = tags
	g.World = world
	g.Turn = PlayerTurn
	g.TurnCounter = 0
	g.DebugMode = debugMode
	
	// Initialize FOV for the starting level
	for _, plr := range g.World.Query(g.WorldTags["players"]) {
		pos := plr.Components[position].(*Position)
		g.Map.CurrentLevel.PlayerVisible.Compute(g.Map.CurrentLevel, pos.X, pos.Y, 8)
	}
	
	return g
}

func (g *Game) Layout(w, h int) (int, int) {
	gd := NewGameData()
	return gd.TileWidth * gd.ScreenWidth, gd.TileHeight * gd.ScreenHeight
}

func (g *Game) Update() error {
	g.TurnCounter++
	if g.Turn == PlayerTurn && g.TurnCounter > 20 {
		TakePlayerAction(g)
	}
	if g.Turn == MonsterTurn {
		UpdateMonster(g)
	}

	return nil
}

// Draw is called each draw cycle and is where we will blit.
func (g *Game) Draw(screen *ebiten.Image) {
	//Draw the Map
	level := g.Map.CurrentLevel
	level.DrawLevel(screen, g.DebugMode)
	ProcessRenderables(g, level, screen)
	ProcessUserLog(g, screen)
	ProcessHUD(g, screen)
}

func main() {
	debugMode := flag.Bool("d", false, "Enable debug mode")
	flag.Parse()
	
	g := NewGame(*debugMode)
	ebiten.SetFullscreen(true)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	
	title := "GoRog"
	if g.DebugMode {
		title += " (Debug Mode)"
	}
	ebiten.SetWindowTitle(title)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
