package main

import (
	"image"

	"github.com/hajimehoshi/ebiten"
)

const (
	north = iota
	east
	south
	west
)

const (
	dead = iota
	idle
	attack
	walk
	cast
	activate
)

const (
	hostile = iota - 1
	neutral
	friendly
)

type actor struct {

	// actor
	x       int
	y       int
	name    string
	coord   vec64
	frame   int
	state   int
	faction int
	// movespeed float64
	effects   map[string]struct{}
	direction int
	anims     map[int][]*ebiten.Image
}

// EFFECTS NOTES: you can get a pseudo set with the following:
// map[string] struct{}
// value, ok := yourmap[key]
// struct{}{}

func spawn_actor(x int, y int, name string, anims map[int][]*ebiten.Image) *actor {
	var a = actor{x: x, y: y}
	a.name = name
	a.anims = anims
	a.frame = 0
	a.direction = 0

	a.coord.x, a.coord.y = cartesianToIso(float64(a.x), float64(a.y))
	a.effects = map[string]struct{}{}
	a.state = idle
	return &a
}

func generateActorSprites(p *ebiten.Image, num_rows int, size int) map[int][]*ebiten.Image {
	anim := make(map[int][]*ebiten.Image)
	num_frames := 10
	directions := 4
	for y := 0; y < num_rows; y++ {
		for x := 0; x < num_frames*directions; x++ {
			anim[4-y] = append(anim[4-y], p.SubImage(image.Rect(size*x, size*y, size*(x+1), size*(y+1))).(*ebiten.Image))
		}
	}
	return anim
}

func generateCharacterSprites(p *ebiten.Image, size int) map[int][]*ebiten.Image {
	anim := make(map[int][]*ebiten.Image)
	num_poses := 6
	num_angles := 8
	num_frames := 10
	for i := 0; i < num_poses; i++ {
		for a := 0; a < num_angles; a++ {
			for f := 0; f < num_frames; f++ {
				y_offset := i*size*4 + (a/2)*size
				x_offset := (a%2)*num_frames*size + f*size
				anim[i] = append(anim[i], p.SubImage(image.Rect(x_offset, y_offset, x_offset+size, y_offset+size)).(*ebiten.Image))
			}
		}
	}
	return anim
}

func wayfind(x1 int, y1 int, x2 int, y2 int) int {
	d := 0
	xy_diff := vec{x: x1 - x2, y: y1 - y2}

	switch {
	case xy_diff.x == 1 && xy_diff.y == 0:
		d = 4
	case xy_diff.x == -1 && xy_diff.y == 0:
		d = 0
	case xy_diff.x == 0 && xy_diff.y == 1:
		d = 6
	case xy_diff.x == 0 && xy_diff.y == -1:
		d = 2
	// UP
	case xy_diff.x == 1 && xy_diff.y == 1:
		d = 5
	case xy_diff.x == -1 && xy_diff.y == 1:
		d = 7
	case xy_diff.x == 1 && xy_diff.y == -1:
		d = 3
	case xy_diff.x == -1 && xy_diff.y == -1:
		d = 1
	}
	return d
}
