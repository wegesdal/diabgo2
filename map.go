package main

import (
	"image"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
)

const (
	sentinal = iota
	grass_tile
	dirt_tile
	road_tile
	cobble_tile
	sand_tile
	river_tile
	block_tile1
	block_tile2
	parlor_white
	parlor_black
)

const (
	chunkSize = 8
)

func generateDoodads(tilesImage *ebiten.Image) []*ebiten.Image {
	var doodads []*ebiten.Image
	tileSize := 512
	for j := 0; j < 4; j++ {
		for i := 0; i < 8; i++ {
			doodads = append(doodads, doodadsImage.SubImage(image.Rect(i*tileSize, j*tileSize, (i+1)*tileSize, (j+1)*tileSize)).(*ebiten.Image))
		}
	}
	return doodads
}

// TODO: Clean up art file and import structure
func generateTiles(tilesImage *ebiten.Image) []*ebiten.Image {
	var tiles []*ebiten.Image

	for i := 0; i < 22; i++ {
		sx := (i % 5) * tileSize
		sy := (i / 5) * tileSize
		tiles = append(tiles, tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image))
	}
	return tiles
}

func generateMap() [3][3][chunkSize][chunkSize]*node {
	var levelData = [3][3][chunkSize][chunkSize]*node{}
	for cx := 0; cx < 3; cx++ {
		for cy := 0; cy < 3; cy++ {
			for x := 0; x < chunkSize; x++ {
				for y := 0; y < chunkSize; y++ {
					levelData[cx][cy][x][y] = &node{x: x + cx*chunkSize, y: y + cy*chunkSize, tile: sentinal}
					levelData[cx][cy][x][y].walkable = true
				}
			}
		}
	}
	return levelData
}

func flattenMap() [chunkSize * 3][chunkSize * 3]*node {
	var flatMap = [chunkSize * 3][chunkSize * 3]*node{}
	for cx := 0; cx < 3; cx++ {
		for cy := 0; cy < 3; cy++ {
			for x := 0; x < chunkSize; x++ {
				for y := 0; y < chunkSize; y++ {
					flatMap[x+cx*chunkSize][y+cy*chunkSize] = levelData[cx][cy][x][y]
				}
			}
		}
	}
	return flatMap
}

func compute_noise(x int, y int) (int, bool, bool) {
	var t int
	var w bool
	var v bool
	noise := perlin(float64(x), float64(y), gradient)

	if noise > 40000.0 {
		t = 17
		w = false
		v = true
	} else {
		if (x%2+y%2)%2 == 0 {
			t = 20
		} else {
			if rand.Intn(1000) == 0 {
				spawnBoss(x, y)
			}
			t = 1
		}
		w = true
		v = false
	}

	return t, w, v
}

func spawnCreep(x int, y int) {
	creepActor := spawn_actor(x, y, "creep", creepAnim)
	c := spawn_character(creepActor)
	c.actor.faction = hostile
	c.prange = 100000.0
	c.arange = 5000.0
	actors = append(actors, creepActor)
	characters = append(characters, c)
}

func spawnBoss(x int, y int) {
	bossActor := spawn_actor(x, y, "boss", bossAnim)
	c := spawn_character(bossActor)
	c.actor.faction = hostile
	c.prange = 100000.0
	c.arange = 5000.0
	actors = append(actors, bossActor)
	characters = append(characters, c)
}
