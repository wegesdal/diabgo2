package main

import (
	"image"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
)

func wall_gen(x int, y int, levelData [32][32]*node) [32][32]*node {

	blocks := 6

	for blocks > 0 {
		levelData[x][y].tile = 4
		levelData[x][y].walkable = false

		if x < 30 && x > 0 && y < 30 && y > 0 {
			d6 := rand.Intn(6)
			if d6 < 2 {
				x++
			} else if d6 < 4 {
				y++
			}
		}
		blocks--
	}
	return levelData
}

// TODO: Clean up art file and import structure
func generateTiles(tilesImage *ebiten.Image) []*ebiten.Image {
	var tiles []*ebiten.Image
	tileSize := 64

	// ground
	sx := 128
	sy := 640
	tiles = append(tiles, tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image))

	// road
	sx = 128
	sy = 256
	tiles = append(tiles, tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image))

	// bridge
	sx = 128
	sy = 320
	tiles = append(tiles, tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image))

	// river
	sx = 896
	sy = 64
	tiles = append(tiles, tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image))

	// block
	sx = 128
	sy = 576
	tiles = append(tiles, tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image))

	return tiles
}

func generateMap() ([32][32]*node, *node) {
	endOfTheRoad := &node{x: rand.Intn(31), y: 31}
	var levelData = [32][32]*node{}

	for x := 0; x < 32; x++ {
		for y := 0; y < 32; y++ {
			levelData[x][y] = &node{x: x, y: y, tile: 0}
			if levelData[x][y].tile < 2 {
				levelData[x][y].walkable = true
			} else {
				levelData[x][y].walkable = false
			}
		}
	}

	// make some walls
	for i := 0; i < 10; i++ {
		start_x := rand.Intn(25) + 4
		start_y := rand.Intn(25) + 4
		levelData = wall_gen(start_x, start_y, levelData)
	}

	// generate a path
	road_start := &node{x: rand.Intn(31), y: 0}

	road := Astar(road_start, endOfTheRoad, levelData)

	// generate a river
	river_start := &node{x: 0, y: rand.Intn(31)}
	river := Astar(river_start, &node{x: 31, y: rand.Intn(31)}, levelData)

	// bake the road onto the array
	for _, node := range road {
		levelData[node.x][node.y].tile = 1
		levelData[node.x+1][node.y].tile = 1
		levelData[node.x][node.y].walkable = true
		levelData[node.x+1][node.y].walkable = true
	}
	// bake the river onto the array
	river = append(river, river_start)

	for _, node := range river {
		if levelData[node.x][node.y].tile == 1 {
			levelData[node.x][node.y].tile = 2
			levelData[node.x][node.y+1].tile = 2
			levelData[node.x][node.y+2].tile = 2
		} else {
			levelData[node.x][node.y].tile = 3
			levelData[node.x][node.y].walkable = false
			levelData[node.x][node.y+1].tile = 3
			levelData[node.x][node.y+1].walkable = false
		}
	}

	return levelData, endOfTheRoad
}
