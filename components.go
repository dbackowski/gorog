package main

import "github.com/hajimehoshi/ebiten/v2"

type Player struct{}

type Monster struct {
	Name string
}

type Position struct {
	X int
	Y int
}

type Renderable struct {
	Image *ebiten.Image
}

type Movable struct{}
