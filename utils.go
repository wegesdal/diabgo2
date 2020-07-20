package main

import (
	"math"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"golang.org/x/image/colornames"
)

func sub_vec64(a vec64, b vec64) vec64 {
	return vec64{x: a.x - b.x, y: a.y - b.y}
}

func cartesianToIso(x, y float64) (float64, float64) {
	rx := (x - y) * float64(tileSize/2)
	ry := (x + y) * float64(tileSize/4)
	return rx, ry
}

func isoToCartesian(x, y float64) (float64, float64) {
	rx := (x/float64(tileSize/2) + y/float64(tileSize/4)) / 2
	ry := (y/float64(tileSize/4) - (x / float64(tileSize/2))) / 2
	return rx, ry
}

// adapted from C http://web.archive.org/web/20110314030147/http://paulbourke.net/geometry/insidepoly/

func insidePolygon(polygon []vec64, N int, p vec64) bool {
	counter := 0
	var i int
	var xinters float64
	var p1, p2 vec64

	p1 = polygon[0]
	for i = 1; i <= N; i++ {
		p2 = polygon[i%N]
		if p.y > math.Min(p1.y, p2.y) {
			if p.y <= math.Max(p1.y, p2.y) {
				if p.x <= math.Max(p1.x, p2.x) {
					if p1.y != p2.y {
						xinters = (p.y-p1.y)*(p2.x-p1.x)/(p2.y-p1.y) + p1.x
						if p1.x == p2.x || p.x <= xinters {

							counter++
						}
					}
				}
			}
		}
		p1 = p2
	}
	var b bool
	if counter%2 == 0 {
		b = false
	} else {
		b = true
	}
	return b
}

func isoSquare(g *Game, screen *ebiten.Image, centerXY vec64, size int, faction int) {
	// 	imd.Color = factionColor(faction, light)
	hs := float64(size / 2)
	// y_offset := -10.0
	// 	centerXY = pixel.Vec.Add(centerXY, pixel.Vec{X: 0, Y: y_offset})

	v1x, v1y := cartesianToIso(-hs, hs-1)
	v2x, v2y := cartesianToIso(hs, hs-1)
	v3x, v3y := cartesianToIso(hs, -hs-1)
	v4x, v4y := cartesianToIso(-hs, -hs-1)

	cx, cy := centerXY.x-g.CamPosX+float64(g.windowWidth/2.0), centerXY.y+g.CamPosY+float64(g.windowHeight/2.0)+40.0

	ebitenutil.DrawLine(screen, v1x+cx, v1y+cy, v2x+cx, v2y+cy, factionColor(faction, light))
	ebitenutil.DrawLine(screen, v2x+cx, v2y+cy, v3x+cx, v3y+cy, factionColor(faction, light))
	ebitenutil.DrawLine(screen, v3x+cx, v3y+cy, v4x+cx, v4y+cy, factionColor(faction, light))
	ebitenutil.DrawLine(screen, v4x+cx, v4y+cy, v1x+cx, v1y+cy, factionColor(faction, light))
}

func isoTargetDebug(g *Game, screen *ebiten.Image, coord vec64) {

	offset := vec64{x: 32.0, y: -30.0}

	cx, cy := coord.x-g.CamPosX+float64(g.windowWidth/2.0)+offset.x, coord.y+g.CamPosY+float64(g.windowHeight/2.0)+offset.y

	ebitenutil.DrawLine(screen, cx+64.0, cy-64.0, cx+64.0, cy+64.0, colornames.Pink)
	ebitenutil.DrawLine(screen, cx-64.0, cy-64.0, cx-64.0, cy+64.0, colornames.Pink)
}

func findOpenNode(levelData [chunkSize][chunkSize]*node) *node {
	// x := rand.Intn(31)
	// y := rand.Intn(31)
	// for !levelData[x][y].walkable {
	// 	x = rand.Intn(31)
	// 	y = rand.Intn(31)
	// }
	return &node{x: 0, y: 0}
}

func getTileXY(g *Game) (int, int) {
	//Get the cursor position
	mx, my := ebiten.CursorPosition()
	//Offset for center
	fmx := float64(mx) - float64(g.windowWidth)/2.0
	fmy := float64(my) - float64(g.windowHeight)/2.0

	x, y := fmx+g.CamPosX, fmy-g.CamPosY

	//Do a half tile mouse shift because of our perspective
	x -= .5 * float64(g.tileSize)
	y -= .5 * float64(g.tileSize)
	//Convert isometric
	imx, imy := isoToCartesian(x, y)
	return int(imx) + ((1 - gx) * chunkSize), int(imy+0.5) + ((1 - gy) * chunkSize)
}
