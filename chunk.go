package main

func chunk() {
	if player.actor.x < chunkSize*gx {
		levelData[2] = levelData[1]
		levelData[1] = levelData[0]

		newRow := [3][chunkSize][chunkSize]*node{}
		for cy := 0; cy < 3; cy++ {
			for x := 0; x < chunkSize; x++ {
				for y := 0; y < chunkSize; y++ {
					newRow[cy][x][y] = &node{x: x + chunkSize*(gx-2), y: y + cy*chunkSize - (1-gy)*chunkSize, tile: sentinal}
					newRow[cy][x][y].blocks_vision = true
				}
			}
		}

		levelData[0] = newRow

		gx--

		flatMap = flattenMap()
	}

	if player.actor.x > chunkSize*(gx+1) {
		levelData[0] = levelData[1]
		levelData[1] = levelData[2]

		newRow := [3][chunkSize][chunkSize]*node{}
		for cy := 0; cy < 3; cy++ {
			for x := 0; x < chunkSize; x++ {
				for y := 0; y < chunkSize; y++ {
					newRow[cy][x][y] = &node{x: x + chunkSize*(gx+2), y: y + cy*chunkSize - (1-gy)*chunkSize, tile: sentinal}
					newRow[cy][x][y].blocks_vision = true
				}
			}
		}
		levelData[2] = newRow
		gx++
		flatMap = flattenMap()
	}

	if player.actor.y < chunkSize*gy {
		for i := 0; i < 3; i++ {
			levelData[i][2] = levelData[i][1]
			levelData[i][1] = levelData[i][0]
		}
		for i := 0; i < 3; i++ {
			newCol := [chunkSize][chunkSize]*node{}
			for x := 0; x < chunkSize; x++ {
				for y := 0; y < chunkSize; y++ {
					newCol[x][y] = &node{x: x + i*chunkSize - (1-gx)*chunkSize, y: y + chunkSize*(gy-2), tile: sentinal}
					newCol[x][y].blocks_vision = true
				}
			}

			levelData[i][0] = newCol
		}

		gy--
		flatMap = flattenMap()
	}

	if player.actor.y > chunkSize*(gy+1) {
		for i := 0; i < 3; i++ {
			levelData[i][0] = levelData[i][1]
			levelData[i][1] = levelData[i][2]
		}
		for i := 0; i < 3; i++ {
			newCol := [chunkSize][chunkSize]*node{}
			for x := 0; x < chunkSize; x++ {
				for y := 0; y < chunkSize; y++ {
					newCol[x][y] = &node{x: x + i*chunkSize - (1-gx)*chunkSize, y: y + chunkSize*(gy+2), tile: sentinal}
					newCol[x][y].blocks_vision = true
				}
			}

			levelData[i][2] = newCol
		}
		gy++
		flatMap = flattenMap()
	}
}
