package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type vec64 struct {
	x float64
	y float64
}

const tileSize = 64.0

var (
	tilesImage   *ebiten.Image
	levelData    [32][32]*node
	endOfTheRoad *node
	tiles        []*ebiten.Image
)

const (
	windowWidth  = 1280
	windowHeight = 900
	screenWidth  = 640
	screenHeight = 480
)

var (
	playerSheet *ebiten.Image
	player      *character
	actors      []*actor
	characters  []*character
	frames      int
	second      = time.Tick(time.Second)
)

type Game struct {
	Name          string //Name of the game ("gollercoaster for now")
	windowWidth   int
	windowHeight  int
	tileSize      int
	CamPosX       float64
	CamPosY       float64
	CamSpeed      float64
	CamZoom       float64
	CamZoomSpeed  float64
	op            *ebiten.DrawImageOptions
	buffer        *ebiten.Image
	drawToBuffer  bool
	lastMousePosX int
	lastMousePosY int
	count         int
}

func init() {
	var err error
	tilesImage, _, err = ebitenutil.NewImageFromFile("dawnblocker.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	playerSheet, _, err = ebitenutil.NewImageFromFile("floyd.png", ebiten.FilterDefault)
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

func lerp_64(v0x float64, v0y float64, v1x float64, v1y float64, t float64) (float64, float64) {
	return (1-t)*v0x + t*v1x, (1-t)*v0y + t*v1y
}

func (g *Game) Update(screen *ebiten.Image) error {

	clearVisibility(levelData)
	compute_fov(vec{x: player.actor.x, y: player.actor.y}, levelData)

	for _, c := range characters {
		cx, cy := cartesianToIso(float64(c.actor.x), float64(c.actor.y))
		c.actor.coord.x, c.actor.coord.y = lerp_64(c.actor.coord.x, c.actor.coord.y, cx, cy, 0.06)
	}

	g.CamPosX, g.CamPosY = lerp_64(g.CamPosX, g.CamPosY, player.actor.coord.x, -player.actor.coord.y, 0.03)

	dt := 2.0 / 60
	// Write your game's logical update.

	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		g.CamPosX -= g.CamSpeed * dt / g.CamZoom
		g.drawToBuffer = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.CamPosX += g.CamSpeed * dt / g.CamZoom
		g.drawToBuffer = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		g.CamPosY -= g.CamSpeed * dt / g.CamZoom
		g.drawToBuffer = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		g.CamPosY += g.CamSpeed * dt / g.CamZoom
		g.drawToBuffer = true
	}
	_, sY := ebiten.Wheel()
	g.CamZoom *= math.Pow(g.CamZoomSpeed, sY)

	if sY != 0 {
		g.drawToBuffer = true
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		// fmt.Printf("x: %d, y: %d", mx, my)

		player.actor.state = walk
		player.target = player

		//Get the cursor position
		mx, my := ebiten.CursorPosition()
		//Offset for center
		fmx := float64(mx) - float64(g.windowWidth)/2.0
		fmy := float64(my) - float64(g.windowHeight)/2.0
		// x, y := float64(mx)+float64(g.windowWidth/2.0), float64(my)+float64(g.windowHeight/2.0)
		//Translate it to game coordinates
		x, y := (float64(fmx/g.CamZoom) + g.CamPosX), float64(fmy/g.CamZoom)-g.CamPosY

		//Do a half tile mouse shift because of our perspective
		x -= .5 * float64(g.tileSize)
		y -= .5 * float64(g.tileSize)
		//Convert isometric
		imx, imy := isoToCartesian(x, y)

		tileX := int(imx)
		tileY := int(imy)

		// for _, c := range characters {
		// offset y so targeting box is above model
		// y_offset := 80.0
		// diff := pixel.Vec.Add(pixel.Vec.Sub(c.actor.coord, cam.Unproject(win.MousePosition())), pixel.Vec{X: 0, Y: y_offset})
		// if math.Abs(diff.X) < 50 && math.Abs(diff.Y) < 100 && c != player {
		// 	player.target = c
		// 	break
		// }
		// }
		if tileX < len(levelData[0]) && tileY < len(levelData[0]) && tileX >= 0 && tileY >= 0 {
			player.dest = &node{x: tileX, y: tileY}
		}
	}

	if g.count%3 == 0 {
		characterStateMachine(characters, levelData)

		// terminalStateMachine(actors)
		actors, characters = removeDeadCharacters(actors, characters)

	}
	g.count++
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.buffer.Clear()
	screen.Fill(color.RGBA{0x10, 0x10, 0x10, 1})

	for x := 0; x < len(levelData[0]); x++ {
		for y := 0; y < len(levelData[0]); y++ {
			xi, yi := cartesianToIso(float64(x), float64(y))
			g.op.GeoM.Reset()
			//Translate for isometric
			g.op.GeoM.Translate(float64(xi), float64(yi))
			//Translate for camera position
			g.op.GeoM.Translate(-g.CamPosX, g.CamPosY)
			//Scale for camera zoom
			g.op.GeoM.Scale(g.CamZoom, g.CamZoom)
			//Translate for center of screen offset
			g.op.GeoM.Translate(float64(g.windowWidth/2.0), float64(g.windowHeight/2.0))
			if levelData[x][y].visible {
				t := tiles[levelData[x][y].tile]
				screen.DrawImage(t, g.op)
			}
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
			g.op.GeoM.Reset()
			//Translate for isometric
			g.op.GeoM.Translate(float64(a.coord.x), float64(a.coord.y))
			g.op.GeoM.Translate(-112.0, -96.0)

			//Translate for camera position
			g.op.GeoM.Translate(-g.CamPosX, g.CamPosY)
			//Scale for camera zoom
			g.op.GeoM.Scale(g.CamZoom, g.CamZoom)
			g.op.GeoM.Translate(float64(g.windowWidth/2.0), float64(g.windowHeight/2.0))

			// The screen should be avoided as a render source
			// If I want the tiles to overlap the feet of the gopher, I'll need to
			// Create another render source for the gopher to prevent conflicting render calls
			// And then insert it into the painter's algorithm
			screen.DrawImage(a.anims[a.state][(a.frame+startingFrame)], g.op)

			// targetRect(player.target.actor.coord, imd, player.target.actor.faction)
			// isoSquare(player.target.actor.coord, 3, imd, player.target.actor.faction)

		} else {
			// DRAW WIDGETS
			widget_coord := sub_vec64(a.coord, vec64{x: 0.0, y: -60.0})

			g.op.GeoM.Reset()
			//Translate for isometric
			g.op.GeoM.Translate(float64(widget_coord.x), float64(widget_coord.y))
			//Translate for camera position
			g.op.GeoM.Translate(-g.CamPosX, g.CamPosY)
			//Scale for camera zoom
			g.op.GeoM.Scale(g.CamZoom, g.CamZoom)

			// raise y by 60 after i make a vec add fn
			screen.DrawImage(a.anims[4][(a.frame+startingFrame)], g.op)

			// isoSquare(a.coord, 3, imd, neutral)
		}
	}
	drawHealthPlates(g, screen, characters)

	//Get the cursor position
	mx, my := ebiten.CursorPosition()
	//Offset for center
	fmx := float64(mx) - float64(g.windowWidth)/2.0
	fmy := float64(my) - float64(g.windowHeight)/2.0
	// x, y := float64(mx)+float64(g.windowWidth/2.0), float64(my)+float64(g.windowHeight/2.0)
	//Translate it to game coordinates
	x, y := (float64(fmx/g.CamZoom) + g.CamPosX), float64(fmy/g.CamZoom)-g.CamPosY

	//Do a half tile mouse shift because of our perspective
	x -= .5 * float64(g.tileSize)
	y -= .5 * float64(g.tileSize)
	//Convert isometric
	// imx, imy := isoToCartesian(x, y)

	// tileX := int(imx)
	// tileY := int(imy)
	ebitenutil.DebugPrint(
		screen,
		fmt.Sprintf("TPS: %0.2f\n", ebiten.CurrentTPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return windowWidth, windowHeight
}

func main() {
	g := &Game{
		Name:         "Diabgo",
		windowWidth:  1280,
		windowHeight: 720,
		tileSize:     64,
		CamPosX:      0,
		CamPosY:      0,
		CamSpeed:     500,
		CamZoom:      1,
		CamZoomSpeed: 1.2,
		op:           &ebiten.DrawImageOptions{},
		drawToBuffer: true,
	}

	g.buffer, _ = ebiten.NewImage(g.windowWidth, g.windowHeight, ebiten.FilterDefault)

	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Diabgo")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}

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
