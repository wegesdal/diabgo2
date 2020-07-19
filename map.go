package main

import (
	"image"

	"github.com/hajimehoshi/ebiten"
)

const (
	grass_tile = iota
	dirt_tile
	road_tile
	cobble_tile
	bridge_tile
	river_tile
	block_tile1
	block_tile2
)

const (
	mapSize = 128
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

	// bridge
	sx = 128
	sy = 512
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

	return tiles
}

func generateMap() [2][mapSize][mapSize]*node {
	// endOfTheRoad := &node{x: 0, y: rand.Intn(mapSize - 1)}
	var levelData = [2][mapSize][mapSize]*node{}

	// gradient := generateGradient()

	// var road []*node
	// f := 1.0

	// for len(road) == 0 {

	for x := 0; x < mapSize; x++ {
		for y := 0; y < mapSize; y++ {
			levelData[1][x][y] = &node{x: x, y: y, tile: 0}
			// noise := perlin(float64(x), float64(y), gradient)
			// if noise > 30000.0*f {
			// 	levelData[0][x][y] = &node{x: x, y: y, tile: river_tile}
			// } else if noise > 10000.0*f {
			// 	levelData[0][x][y] = &node{x: x, y: y, tile: dirt_tile}
			// } else if noise > -2000.0*f {
			// 	levelData[0][x][y] = &node{x: x, y: y, tile: grass_tile}
			// } else if noise > -10000.0*f {
			// 	levelData[0][x][y] = &node{x: x, y: y, tile: cobble_tile}
			// } else if noise > -20000.0*f {
			// 	levelData[0][x][y] = &node{x: x, y: y, tile: block_tile1}
			// } else {
			// 	levelData[0][x][y] = &node{x: x, y: y, tile: block_tile2}
			// }
			levelData[0][x][y] = &node{x: x, y: y, tile: dirt_tile}
			if levelData[0][x][y].tile < river_tile {
				levelData[0][x][y].walkable = true
			} else {
				levelData[0][x][y].walkable = false
			}
		}
	}

	// road_start := &node{x: mapSize - 1, y: rand.Intn(mapSize - 1)}
	// levelData[0][road_start.x][road_start.y].tile = road_tile
	// road = Astar(road_start, endOfTheRoad, levelData[0], false)
	// // 	f++
	// // }

	// // bake the road onto the array
	// for _, node := range road {
	// 	levelData[0][node.x][node.y].tile = road_tile
	// 	levelData[0][node.x][node.y].walkable = true
	// 	if node.x+1 < mapSize-1 {
	// 		levelData[0][node.x+1][node.y].tile = road_tile
	// 		levelData[0][node.x+1][node.y].walkable = true
	// 	}
	// }

	// generate a river
	// river_start := &node{x: 0, y: rand.Intn(mapSize - 1)}
	// river := Astar(river_start, &node{x: mapSize - 1, y: rand.Intn(mapSize - 1)}, levelData[0], false)

	// // bake the river onto the array
	// river = append(river, river_start)

	// for _, node := range river {
	// 	if levelData[0][node.x][node.y].tile == road_tile {
	// 		levelData[0][node.x][node.y].tile = bridge_tile
	// 		if node.y+1 < mapSize-1 {
	// 			levelData[0][node.x][node.y+1].tile = bridge_tile
	// 		}
	// 		if node.y+2 < mapSize-1 {
	// 			levelData[0][node.x][node.y+2].tile = bridge_tile
	// 		}
	// 	} else {
	// 		levelData[0][node.x][node.y].tile = river_tile
	// 		levelData[0][node.x][node.y].walkable = false
	// 		if node.y+1 < mapSize-1 {
	// 			levelData[0][node.x][node.y+1].tile = river_tile
	// 			levelData[0][node.x][node.y+1].walkable = false
	// 		}
	// 	}
	// }

	return levelData
}
