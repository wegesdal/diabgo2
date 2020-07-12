package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type vec64 struct {
	x float64
	y float64
}

var tileSize = 64.0

var tilesImage *ebiten.Image

var levelData [32][32]*node
var endOfTheRoad *node
var tiles []*ebiten.Image

func init() {
	var err error
	tilesImage, _, err = ebitenutil.NewImageFromFile("dawnblocker.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	levelData, endOfTheRoad = generateMap()
	tiles = generateTiles(tilesImage)
}

type Game struct{}

func (g *Game) Update(screen *ebiten.Image) error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x10, 0x10, 0x10, 1})
	ebitenutil.DebugPrint(screen, "Hello, Will!")

	for x := 0; x < len(levelData[0]); x++ {
		for y := 0; y < len(levelData[0]); y++ {
			isoCoords := cartesianToIso(vec64{x: float64(x), y: float64(y)})
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(isoCoords.x, isoCoords.y)
			screen.DrawImage(tiles[levelData[x][y].tile], op)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1280, 900
}

func main() {
	ebiten.SetWindowSize(1280, 900)
	ebiten.SetWindowTitle("Diabgo")

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

func cartesianToIso(pt vec64) vec64 {
	return vec64{x: (pt.x - pt.y) * (tileSize / 2), y: (pt.x + pt.y) * (tileSize / 4)}
}

func isoToCartesian(pt vec64) vec64 {
	x := pt.x*(2.0/tileSize) + pt.y*(4/tileSize)
	y := ((pt.y * 4.0 / tileSize) - x) / 2
	return vec64{x: x + y, y: y}
}
