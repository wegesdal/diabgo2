package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type vec64 struct {
	x float64
	y float64
}

var (
	tileSize     = 64.0
	tilesImage   *ebiten.Image
	levelData    [32][32]*node
	endOfTheRoad *node
	tiles        []*ebiten.Image
)

var (
	playerSheet *ebiten.Image
	player      *character
	actors      []*actor
	characters  []*character
	frames      int
	second      = time.Tick(time.Second)
)

func init() {
	var err error
	tilesImage, _, err = ebitenutil.NewImageFromFile("dawnblocker.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	playerSheet, _, err = ebitenutil.NewImageFromFile("gopher8.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	levelData, endOfTheRoad = generateMap()
	tiles = generateTiles(tilesImage)

	player_anim := generateCharacterSprites(playerSheet, 256)
	player_spawn := findOpenNode(levelData)
	act := spawn_actor(player_spawn.x, player_spawn.y, "player", player_anim)
	player = spawn_character(act)
	player.maxhp = 40
	player.hp = 40
	player.actor.faction = friendly
	player.prange = 0.0
	player.arange = 5000.0
	actors = append(actors, act)
	characters = append(characters, player)

}

type Game struct {
	count int
}

func (g *Game) Update(screen *ebiten.Image) error {

	characterStateMachine(characters, levelData)
	// terminalStateMachine(actors)
	actors, characters = removeDeadCharacters(actors, characters)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x10, 0x10, 0x10, 1})
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.CurrentFPS()))

	drawHealthPlates(screen, characters)

	for x := 0; x < len(levelData[0]); x++ {
		for y := 0; y < len(levelData[0]); y++ {
			isoCoords := cartesianToIso(vec64{x: float64(x), y: float64(y)})
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(isoCoords.x, isoCoords.y)
			screen.DrawImage(tiles[levelData[x][y].tile], op)
		}
	}

	// DRAW ACTORS
	for _, a := range actors {

		startingFrame := 0
		// half_length := len(a.anims[a.state]) / 2
		// i := isoToCartesian(a.coord)
		// draw actors
		// offset := 0.2
		// if x == int(i.x+offset) && y == int(i.y+offset) {
		// DRAW CHARACTER
		// the length of anims tells you if this is a character or item
		// characters will have an anims length of 6
		// widgets will have an anims length of 1
		startingFrame = a.direction * 10

		if len(a.anims) == 6 {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(0.0, 0.0)

			// The screen should be avoided as a render source
			// If I want the tiles to overlap the feet of the gopher, I'll need to
			// Create another render source for the gopher to prevent conflicting render calls
			// And then insert it into the painter's algorithm
			screen.DrawImage(a.anims[a.state][(a.frame+startingFrame)], op)

			// targetRect(player.target.actor.coord, imd, player.target.actor.faction)
			// isoSquare(player.target.actor.coord, 3, imd, player.target.actor.faction)

		} else {
			// DRAW WIDGETS
			widget_coord := sub_vec64(a.coord, vec64{x: 0.0, y: -60.0})
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(widget_coord.x, widget_coord.y)

			// raise y by 60 after i make a vec add fn
			screen.DrawImage(a.anims[4][(a.frame+startingFrame)], op)

			// isoSquare(a.coord, 3, imd, neutral)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 800, 600
}

func main() {
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Diabgo")

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

func cartesianToIso(pt vec64) vec64 {
	return vec64{x: (pt.x - pt.y) * (tileSize / 2.0), y: (pt.x + pt.y) * (tileSize / 4.0)}
}

func isoToCartesian(pt vec64) vec64 {
	x := pt.x*(2.0/tileSize) + pt.y*(4.0/tileSize)
	y := ((pt.y * 4.0 / tileSize) - x) / 2.0
	return vec64{x: x + y, y: y}
}

// func isoSquare(centerXY vec64, size int, imd *imdraw.IMDraw, faction int) {
// 	imd.Color = factionColor(faction, light)
// 	hs := float64(size / 2)
// 	y_offset := -10.0
// 	centerXY = pixel.Vec.Add(centerXY, pixel.Vec{X: 0, Y: y_offset})
// 	imd.Push(pixel.Vec.Add(centerXY, cartesianToIso(pixel.Vec{X: -hs, Y: -hs})))
// 	imd.Push(pixel.Vec.Add(centerXY, cartesianToIso(pixel.Vec{X: -hs, Y: hs})))
// 	imd.Push(pixel.Vec.Add(centerXY, cartesianToIso(pixel.Vec{X: hs, Y: hs})))
// 	imd.Push(pixel.Vec.Add(centerXY, cartesianToIso(pixel.Vec{X: hs, Y: -hs})))
// 	imd.Polygon(1)
// }

func findOpenNode(levelData [32][32]*node) *node {
	x := rand.Intn(31)
	y := rand.Intn(31)
	for levelData[x][y].tile != 1 {
		x = rand.Intn(31)
		y = rand.Intn(31)
	}
	return &node{x: x, y: y}
}
