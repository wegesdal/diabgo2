package main

import "math/rand"

func lerp_f(a0 float64, a1 float64, w float64) float64 {
	return (1.0-w)*a0 + w*a1
}

// Computes the dot product of the distance and gradient vectors.
func dotGridGradient(ix int, iy int, x float64, y float64, gradient [128][128][2]float64) float64 {

	// Precomputed (or otherwise) gradient vectors at each grid node
	// Compute the distance vector
	dx := x - float64(ix)
	dy := y - float64(iy)

	// Compute the dot-product
	return (dx*gradient[iy][ix][0] + dy*gradient[iy][ix][1])
}

// Compute Perlin noise at coordinates x, y
func perlin(x float64, y float64, gradient [128][128][2]float64) float64 {

	// Determine grid cell coordinates
	x0 := int(127 * x / (chunkSize * 128))
	if x0 < 0 {
		x0 += 31
	}
	x1 := x0 + 1
	y0 := int(127 * y / (chunkSize * 128))
	if y0 < 0 {
		y0 += 31
	}
	y1 := y0 + 1

	// Determine interpolation weights
	// Could also use higher order polynomial/s-curve here
	sx := x - float64(x0)
	sy := y - float64(y0)

	// Interpolate between grid point gradients
	var (
		n0    float64
		n1    float64
		ix0   float64
		ix1   float64
		value float64
	)

	n0 = dotGridGradient(x0, y0, x, y, gradient)
	n1 = dotGridGradient(x1, y0, x, y, gradient)
	ix0 = lerp_f(n0, n1, sx)

	n0 = dotGridGradient(x0, y1, x, y, gradient)
	n1 = dotGridGradient(x1, y1, x, y, gradient)
	ix1 = lerp_f(n0, n1, sx)

	value = lerp_f(ix0, ix1, sy)
	return value
}

func generateGradient() [128][128][2]float64 {
	var gradient [128][128][2]float64
	for x := 0; x < len(gradient); x++ {
		for y := 0; y < len(gradient[0]); y++ {
			for z := 0; z < len(gradient[0][0]); z++ {
				gradient[x][y][z] = rand.Float64()*2.0 - 1.0
			}
		}
	}
	return gradient
}
