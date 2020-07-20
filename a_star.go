package main

import (
	"math"
)

func add(x int, y int) int {
	return x + y
}

type node struct {
	x             int
	y             int
	parent        *node
	H             int
	G             int
	tile          int
	walkable      bool
	visible       bool
	blocks_vision bool
}

func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

// manhattan distance is for 4-way movement calculations
func manhattan_distance(a node, b node) int {
	return int(math.Abs(float64(a.x-b.x)) + math.Abs(float64(a.y-b.y)))
}

func diagonal_distance(a node, b node) int {
	return Max((a.x-b.x)*(a.x-b.x), (a.y-b.y)*(a.y-b.y))
}

func walkable(n node, grid [chunkSize * 3][chunkSize * 3]*node) bool {
	if n.x+(1-gx)*chunkSize >= 0 && n.y+(1-gy)*chunkSize >= 0 && n.x+(1-gx)*chunkSize < len(grid) && n.y+(1-gy)*chunkSize < len(grid[0]) {
		return grid[n.x+(1-gx)*chunkSize][n.y+(1-gy)*chunkSize].walkable
	} else {
		return false
	}
}

func retrace(n *node) []*node {
	var path = []*node{}
	var current *node = n

	for current.parent != nil {
		path = append(path, current)
		current = current.parent
	}
	return path
}

// returns the index of the lowest gh value in the slice
func min_gh(nodes []*node) int {
	min := 0
	for i, v := range nodes {
		if v.G+v.H < nodes[min].G+nodes[min].H {
			min = i
		}
	}
	return min
}

// returns true if the location of this node does not appear in the slice
func unique(q *node, nodes []*node) bool {
	for _, n := range nodes {
		if n.x == q.x && n.y == q.y {
			return false
		}
	}
	return true
}

// fast swap remove
func remove(s []*node, i int) []*node {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func Astar(start *node, end *node, grid [chunkSize * 3][chunkSize * 3]*node, diagonals bool) []*node {
	var open = []*node{}
	var closed = []*node{}
	open = append(open, start)

	for len(open) > 0 {

		var c = min_gh(open)

		if open[c].x == end.x && open[c].y == end.y {
			var path = retrace(open[c])
			return path
		}

		var neighbors = []*node{
			&node{x: open[c].x + 1, y: open[c].y},
			&node{x: open[c].x - 1, y: open[c].y},
			&node{x: open[c].x, y: open[c].y + 1},
			&node{x: open[c].x, y: open[c].y - 1},
		}

		if diagonals {
			neighbors = append(neighbors, &node{x: open[c].x + 1, y: open[c].y + 1})
			neighbors = append(neighbors, &node{x: open[c].x - 1, y: open[c].y - 1})
			neighbors = append(neighbors, &node{x: open[c].x + 1, y: open[c].y - 1})
			neighbors = append(neighbors, &node{x: open[c].x - 1, y: open[c].y + 1})
		}

		closed = append(closed, open[c])
		open = remove(open, c)

		for _, n := range neighbors {
			if unique(n, closed) && walkable(*n, grid) {
				if !unique(n, open) {
					var new_G = closed[len(closed)-1].G + 1
					if n.G > new_G {
						n.G = new_G
						n.parent = closed[len(closed)-1]
					}
				} else {
					n.G = closed[len(closed)-1].G + 1
					if diagonals {
						n.H = diagonal_distance(*n, *end)
					} else {
						n.H = manhattan_distance(*n, *end)
					}
					n.parent = closed[len(closed)-1]
					open = append(open, n)
				}
			}
		}
	}
	return open
}
