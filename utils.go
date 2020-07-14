package main

func sub_vec64(a vec64, b vec64) vec64 {
	return vec64{x: a.x - b.x, y: a.y - b.y}
}
