package main

import (
	"image"
	"log"

	"github.com/bytearena/ecs"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var position *ecs.Component
var renderable *ecs.Component
var monster *ecs.Component
var health *ecs.Component
var meleeWeapon *ecs.Component
var armor *ecs.Component
var name *ecs.Component

func InitializeWorld(startingLevel Level) (*ecs.Manager, map[string]ecs.Tag) {
	tags := make(map[string]ecs.Tag)
	manager := ecs.NewManager()

	img, _, err := ebitenutil.NewImageFromFile("assets/EverRogueTileset 1.0 Horizontal.png")

	if err != nil {
		log.Fatal(err)
	}

	playerImg := img.SubImage(image.Rect(2800, 0, 2784, 2800)).(*ebiten.Image)  // player image
	monsterImg := img.SubImage(image.Rect(2784, 0, 2768, 2784)).(*ebiten.Image) // monster image

	startingRoom := startingLevel.Rooms[0]
	x, y := startingRoom.Center()

	player := manager.NewComponent()
	position = manager.NewComponent()
	renderable = manager.NewComponent()
	movable := manager.NewComponent()
	monster = manager.NewComponent()
	health = manager.NewComponent()
	meleeWeapon = manager.NewComponent()
	armor = manager.NewComponent()
	name = manager.NewComponent()

	manager.NewEntity().
		AddComponent(player, Player{}).
		AddComponent(renderable, &Renderable{
			Image: playerImg,
		}).
		AddComponent(movable, Movable{}).
		AddComponent(position, &Position{
			X: x,
			Y: y,
		}).
		AddComponent(health, &Health{
			MaxHealth:     30,
			CurrentHealth: 30,
		}).
		AddComponent(meleeWeapon, &MeleeWeapon{
			Name:          "Fist",
			MinimumDamage: 1,
			MaximumDamage: 3,
			ToHitBonus:    2,
		}).
		AddComponent(armor, &Armor{
			Name:       "Burlap Sack",
			Defense:    1,
			ArmorClass: 1,
		}).
		AddComponent(name, &Name{Label: "Player"})

	for _, room := range startingLevel.Rooms {
		if room.X1 != startingRoom.X1 {
			mX, mY := room.Center()
			manager.NewEntity().
				AddComponent(monster, &Monster{}).
				AddComponent(renderable, &Renderable{
					Image: monsterImg,
				}).
				AddComponent(position, &Position{
					X: mX,
					Y: mY,
				}).
				AddComponent(health, &Health{
					MaxHealth:     10,
					CurrentHealth: 10,
				}).
				AddComponent(meleeWeapon, &MeleeWeapon{
					Name:          "Short Sword",
					MinimumDamage: 2,
					MaximumDamage: 6,
					ToHitBonus:    0,
				}).
				AddComponent(armor, &Armor{
					Name:       "Bone",
					Defense:    3,
					ArmorClass: 4,
				}).
				AddComponent(name, &Name{Label: "Skeleton"})
		}
	}

	players := ecs.BuildTag(player, position, health, meleeWeapon, armor, name)
	tags["players"] = players

	renderables := ecs.BuildTag(renderable, position)
	tags["renderables"] = renderables

	monsters := ecs.BuildTag(monster, position, health, meleeWeapon, armor, name)
	tags["monsters"] = monsters

	return manager, tags
}
