package main

type GameMap struct {
	Dungeons     []Dungeon
	CurrentLevel Level
	LevelIndex   int
}

func NewGameMap() GameMap {
	l := NewLevel(0)
	levels := make([]Level, 0)
	levels = append(levels, l)
	d := Dungeon{Name: "default", Levels: levels}
	dungeons := make([]Dungeon, 0)
	dungeons = append(dungeons, d)
	gm := GameMap{
		Dungeons:     dungeons, 
		CurrentLevel: l,
		LevelIndex:   0,
	}

	return gm
}

// GoDownStairs creates a new level or navigates to an existing deeper level
func (gm *GameMap) GoDownStairs() bool {
	// Check if we already have a level below this one
	if gm.LevelIndex+1 < len(gm.Dungeons[0].Levels) {
		// We already have this level, just move to it
		gm.LevelIndex++
		gm.CurrentLevel = gm.Dungeons[0].Levels[gm.LevelIndex]
		return true
	}
	
	// Create a new level
	newLevel := NewLevel(gm.LevelIndex + 1)
	gm.Dungeons[0].Levels = append(gm.Dungeons[0].Levels, newLevel)
	gm.LevelIndex++
	gm.CurrentLevel = newLevel
	return true
}

// GoUpStairs navigates to the level above if possible
func (gm *GameMap) GoUpStairs() bool {
	if gm.LevelIndex > 0 {
		gm.LevelIndex--
		gm.CurrentLevel = gm.Dungeons[0].Levels[gm.LevelIndex]
		return true
	}
	return false
}
