package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	width  = 80
	height = 24
)

type Position struct {
	x, y int
}

type Entity struct {
	pos    Position
	symbol rune
	health int
}

type Game struct {
	player   Entity
	monsters []Entity
	dungeon  [][]rune
}

func generateDungeon() [][]rune {
	dungeon := make([][]rune, height)
	for i := range dungeon {
		dungeon[i] = make([]rune, width)
		for j := range dungeon[i] {
			if i == 0 || i == height-1 || j == 0 || j == width-1 {
				dungeon[i][j] = '#'
			} else {
				dungeon[i][j] = '.'
			}
		}
	}
	return dungeon
}

func initPlayer() Entity {
	return Entity{
		pos:    Position{width / 2, height / 2},
		symbol: '@',
		health: 10,
	}
}

func initMonsters(count int) []Entity {
	monsters := make([]Entity, count)
	for i := range monsters {
		monsters[i] = Entity{
			pos:    Position{rand.Intn(width-2) + 1, rand.Intn(height-2) + 1},
			symbol: 'M',
			health: 3,
		}
	}
	return monsters
}

func initGame() Game {
	rand.Seed(time.Now().UnixNano())
	return Game{
		player:   initPlayer(),
		monsters: initMonsters(5),
		dungeon:  generateDungeon(),
	}
}

func (g *Game) display() {
	fmt.Print("\033[H\033[2J") // Clear screen (ANSI escape code)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if g.player.pos.x == x && g.player.pos.y == y {
				fmt.Print(string(g.player.symbol))
			} else {
				monsterFound := false
				for _, m := range g.monsters {
					if m.pos.x == x && m.pos.y == y {
						fmt.Print(string(m.symbol))
						monsterFound = true
						break
					}
				}
				if !monsterFound {
					fmt.Print(string(g.dungeon[y][x]))
				}
			}
		}
		fmt.Println()
	}
	fmt.Printf("Player Health: %d\n", g.player.health)
}

func (g *Game) movePlayer(dx, dy int) {
	newX := g.player.pos.x + dx
	newY := g.player.pos.y + dy
	if g.dungeon[newY][newX] != '#' {
		g.player.pos.x = newX
		g.player.pos.y = newY
	}
}

func (g *Game) moveMonsters() {
	for i := range g.monsters {
		dx := rand.Intn(3) - 1
		dy := rand.Intn(3) - 1
		newX := g.monsters[i].pos.x + dx
		newY := g.monsters[i].pos.y + dy
		if g.dungeon[newY][newX] != '#' {
			g.monsters[i].pos.x = newX
			g.monsters[i].pos.y = newY
		}
	}
}

func (g *Game) combat() {
	for i, m := range g.monsters {
		if m.pos == g.player.pos {
			g.player.health--
			g.monsters[i].health--
			if g.monsters[i].health <= 0 {
				g.monsters = append(g.monsters[:i], g.monsters[i+1:]...)
				break
			}
		}
	}
}

func main() {
	game := initGame()

	for {
		game.display()

		var input string
		fmt.Scan(&input)

		switch input {
		case "w":
			game.movePlayer(0, -1)
		case "s":
			game.movePlayer(0, 1)
		case "a":
			game.movePlayer(-1, 0)
		case "d":
			game.movePlayer(1, 0)
		case "q":
			fmt.Println("Thanks for playing!")
			return
		}

		game.moveMonsters()
		game.combat()

		if game.player.health <= 0 {
			fmt.Println("Game Over! You died.")
			return
		}

		if len(game.monsters) == 0 {
			fmt.Println("Congratulations! You defeated all monsters!")
			return
		}
	}
}
