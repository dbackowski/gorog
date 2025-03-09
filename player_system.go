package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func TakePlayerAction(g *Game) {
	turnTaken := false
	players := g.WorldTags["players"]
	x := 0
	y := 0
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		y = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		y = 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		x = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		x = 1
	}

	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		turnTaken = true
	}

	// Check for stairs navigation
	stairsAction := false
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		// Try to go down stairs
		stairsAction = tryUseStairs(g, STAIRS_DOWN)
		turnTaken = stairsAction
	}
	if ebiten.IsKeyPressed(ebiten.KeyU) {
		// Try to go up stairs
		stairsAction = tryUseStairs(g, STAIRS_UP)
		turnTaken = stairsAction
	}

	level := g.Map.CurrentLevel

	for _, result := range g.World.Query(players) {
		pos := result.Components[position].(*Position)
		index := level.GetIndexFromXY(pos.X+x, pos.Y+y)

		tile := level.Tiles[index]
		if tile.Blocked != true {
			level.Tiles[level.GetIndexFromXY(pos.X, pos.Y)].Blocked = false

			pos.X += x
			pos.Y += y
			level.Tiles[index].Blocked = true
			level.PlayerVisible.Compute(level, pos.X, pos.Y, 8)
		} else if x != 0 || y != 0 {
			if level.Tiles[index].TileType != WALL {
				//Its a tile with a monster -- Fight it
				monsterPosition := Position{X: pos.X + x, Y: pos.Y + y}
				AttackSystem(g, pos, &monsterPosition)
			}
		}
	}

	if x != 0 || y != 0 || turnTaken {
		g.Turn = GetNextState(g.Turn)
		g.TurnCounter = 0
	}
}

// tryUseStairs attempts to use stairs if the player is standing on them
func tryUseStairs(g *Game, stairType TileType) bool {
	level := g.Map.CurrentLevel

	// Get player position
	var playerPos *Position
	for _, result := range g.World.Query(g.WorldTags["players"]) {
		playerPos = result.Components[position].(*Position)
		break
	}

	if playerPos == nil {
		return false
	}

	// Check if player is on stairs
	index := level.GetIndexFromXY(playerPos.X, playerPos.Y)
	if level.Tiles[index].TileType != stairType {
		return false
	}

	// Use the stairs
	if stairType == STAIRS_DOWN {
		// Go down to next level
		if g.Map.GoDownStairs() {
			// Update player position to stairs up on the new level
			newLevel := g.Map.CurrentLevel
			if newLevel.StairsUp != nil {
				playerPos.X = newLevel.StairsUp.X
				playerPos.Y = newLevel.StairsUp.Y

				// Update FOV for new level
				newLevel.PlayerVisible.Compute(newLevel, playerPos.X, playerPos.Y, 8)
				return true
			}
		}
	} else if stairType == STAIRS_UP {
		// Go up to previous level
		if g.Map.GoUpStairs() {
			// Update player position to stairs down on the previous level
			newLevel := g.Map.CurrentLevel
			if newLevel.StairsDown != nil {
				playerPos.X = newLevel.StairsDown.X
				playerPos.Y = newLevel.StairsDown.Y

				// Update FOV for new level
				newLevel.PlayerVisible.Compute(newLevel, playerPos.X, playerPos.Y, 8)
				return true
			}
		}
	}

	return false
}
