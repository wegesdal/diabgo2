package main

import (
	"math"
)

type vec struct {
	x int
	y int
}

func clearVisibility(grid [chunkSize][chunkSize]*node) {
	for _, row := range grid {
		for _, node := range row {
			node.visible = false
		}
	}
}

// adapted from https://www.albertford.com/shadowcasting/

func compute_fov(origin vec, grid [chunkSize * 3][chunkSize * 3]*node) {

	grid[origin.x][origin.y].visible = true
	for i := 0; i < 4; i++ {
		quadrant := Quadrant{cardinal: i, origin: origin}
		first_row := Row{depth: 1.0, start_slope: -1.0, end_slope: 1.0}
		scan(first_row, grid, &quadrant)
	}
}

func reveal(tile vec, grid [chunkSize * 3][chunkSize * 3]*node, quadrant *Quadrant) {
	q := transform(quadrant, tile)
	if in_bounds(q.x, q.y, grid) {
		grid[q.x][q.y].visible = true
	}
}

func is_wall(tile vec, grid [chunkSize * 3][chunkSize * 3]*node, quadrant *Quadrant) bool {
	var w bool
	if (vec{}) != tile {
		w = false
		q := transform(quadrant, tile)

		if in_bounds(q.x, q.y, grid) {
			w = grid[q.x][q.y].blocks_vision
		} else {
			w = true
		}
	}
	return w
}

func in_bounds(x int, y int, grid [chunkSize * 3][chunkSize * 3]*node) bool {
	if x < len(grid) && y < len(grid[0]) && x >= 0 && y >= 0 {
		return true
	} else {
		return false
	}
}

func is_floor(tile vec, grid [chunkSize * 3][chunkSize * 3]*node, quadrant *Quadrant) bool {
	var f bool
	if (vec{}) != tile {
		f = false
		q := transform(quadrant, tile)
		if in_bounds(q.x, q.y, grid) {
			f = !grid[q.x][q.y].blocks_vision
		}
	}
	return f
}

func scan(row Row, grid [chunkSize * 3][chunkSize * 3]*node, quadrant *Quadrant) {
	var prev_tile = vec{}
	var tiles = generate_tiles(row)
	for _, tile := range tiles {

		if is_wall(tile, grid, quadrant) || is_symmetric(row, tile) {
			reveal(tile, grid, quadrant)
		}
		if is_wall(prev_tile, grid, quadrant) && is_floor(tile, grid, quadrant) {
			row.start_slope = slope(tile)
		}
		if is_floor(prev_tile, grid, quadrant) && is_wall(tile, grid, quadrant) {
			next_row := next(&row)
			next_row.end_slope = slope(tile)
			scan(next_row, grid, quadrant)
		}
		prev_tile = tile

	}

	if is_floor(prev_tile, grid, quadrant) {
		scan(next(&row), grid, quadrant)
	}

}

type Quadrant struct {
	cardinal int
	origin   vec
}

func transform(self *Quadrant, tile vec) vec {
	row, col := tile.x, tile.y
	var v vec
	switch self.cardinal {
	case north:
		v = vec{x: self.origin.x + col, y: self.origin.y - row}
	case south:
		v = vec{x: self.origin.x + col, y: self.origin.y + row}
	case east:
		v = vec{x: self.origin.x + row, y: self.origin.y + col}
	case west:
		v = vec{x: self.origin.x - row, y: self.origin.y + col}
	}
	return v
}

type Row struct {
	depth       float64
	start_slope float64
	end_slope   float64
}

func generate_tiles(self Row) []vec {
	min_col := round_ties_up(self.depth * self.start_slope)
	max_col := round_ties_down(self.depth * self.end_slope)

	var tiles []vec
	for col := min_col; col < max_col+1; col++ {
		tiles = append(tiles, vec{x: int(self.depth), y: col})
	}
	return tiles
}

func next(self *Row) Row {
	return Row{depth: self.depth + 1.0, start_slope: self.start_slope, end_slope: self.end_slope}
}

func slope(tile vec) float64 {
	row_depth, col := tile.x, tile.y
	return float64(2*col-1) / float64(2*row_depth)
}

func is_symmetric(row Row, tile vec) bool {
	col := tile.y
	return float64(col) >= row.depth*row.start_slope && float64(col) <= row.depth*row.end_slope
}

func round_ties_up(n float64) int {
	return int(math.Floor(n + 0.5))
}

func round_ties_down(n float64) int {
	return int(math.Ceil(n - 0.5))
}
