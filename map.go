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
	tileSize := 64

	// grass
	sx := 2432
	sy := 640
	tiles = append(tiles, tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image))

	// dirt
	sx = 128
	sy = 2304
	tiles = append(tiles, tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image))

	// road
	sx = 2208
	sy = 320
	tiles = append(tiles, tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image))

	// cobble
	sx = 128
	sy = 256
	tiles = append(tiles, tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image))

	// sand
	sx = 1664
	sy = 1472
	tiles = append(tiles, tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image))

	// river
	sx = 896
	sy = 64
	tiles = append(tiles, tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image))

	// block
	sx = 128
	sy = 384
	tiles = append(tiles, tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image))

	// block2
	sx = 128
	sy = 320
	tiles = append(tiles, tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image))

	// parlor white

	sx = 128
	sy = 896
	tiles = append(tiles, tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image))

	// parlor black

	sx = 1664
	sy = 1728
	tiles = append(tiles, tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image))
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
		t = block_tile2
		w = false
		v = true
	} else if noise > 35000.0 {
		t = block_tile1
		w = false
		v = true
	} else if noise > 20000.0 {
		if (x%2+y%2)%2 == 0 {
			t = parlor_white
		} else {
			if rand.Intn(1000) == 0 {
				spawnBoss(x, y)
			}
			t = parlor_black
		}
		w = true
		v = false
	} else if noise > 10000.0 {
		t = cobble_tile
		w = true
		v = false
	} else if noise > 0.0 {
		t = grass_tile
		w = true
		v = false
		if rand.Intn(100) == 0 {
			spawnCreep(x, y)
		}

	} else if noise > -30000.0 {
		t = dirt_tile
		w = true
		v = false
	} else if noise > -70000.0 {
		t = sand_tile
		w = true
		v = false
	} else {
		t = river_tile
		w = false
		v = false
	}
	return t, w, v
}

func spawnCreep(x int, y int) {
	creepActor := spawn_actor(x, y, "creep", creepAnim)
	c := spawn_character(creepActor)
	c.actor.faction = hostile
	c.prange = 18000.0
	c.arange = 5000.0
	actors = append(actors, creepActor)
	characters = append(characters, c)
}

func spawnBoss(x int, y int) {
	bossActor := spawn_actor(x, y, "boss", bossAnim)
	c := spawn_character(bossActor)
	c.actor.faction = hostile
	c.prange = 18000.0
	c.arange = 5000.0
	actors = append(actors, bossActor)
	characters = append(characters, c)
}
