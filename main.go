package main

import (
	"fmt"
	"image/color"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"sort"
	"time"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
)

type vec64 struct {
	x float64
	y float64
}

const tileSize = 64.0

var (
	tilesImage   *ebiten.Image
	doodadsImage *ebiten.Image
	levelData    [2][32][32]*node
	endOfTheRoad *node
	tiles        []*ebiten.Image
	doodads      []*ebiten.Image
)

const (
	screenWidth  = 640
	screenHeight = 480
)

var (
	playerSheet   *ebiten.Image
	creepSheet    *ebiten.Image
	bossSheet     *ebiten.Image
	terminalSheet *ebiten.Image
	player        *character
	actors        []*actor
	characters    []*character
	frames        int
	second        = time.Tick(time.Second)
)

var (
	sampleText  = `Spooky Forest`
	exocet_face font.Face
)

type Game struct {
	Name          string
	windowWidth   int
	windowHeight  int
	tileSize      int
	CamPosX       float64
	CamPosY       float64
	CamSpeed      float64
	op            *ebiten.DrawImageOptions
	buffer        *ebiten.Image
	drawToBuffer  bool
	lastMousePosX int
	lastMousePosY int
	count         int
}

func init() {
	rand.Seed(time.Now().UnixNano())

	exocet_ttf, err := ioutil.ReadFile("assets/fonts/ExocetLight_Medium.ttf")

	tt, err := truetype.Parse(exocet_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	exocet_face = truetype.NewFace(tt, &truetype.Options{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})

	// LOAD ASSETS
	tilesImage, _, err = ebitenutil.NewImageFromFile("assets/sprites/dawnblocker.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	playerSheet, _, err = ebitenutil.NewImageFromFile("assets/sprites/fs_gopher.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	creepSheet, _, err = ebitenutil.NewImageFromFile("assets/sprites/fs_creep.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	bossSheet, _, err = ebitenutil.NewImageFromFile("assets/sprites/fs_diabgopher.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	terminalSheet, _, err = ebitenutil.NewImageFromFile("assets/sprites/fs_terminal.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	doodadsImage, _, err = ebitenutil.NewImageFromFile("assets/sprites/fs_doodads.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	//INITIALIZE WORLD OBJECTS

	levelData, endOfTheRoad = generateMap()
	tiles = generateTiles(tilesImage)
	doodads = generateDoodads(doodadsImage)

	playerAnim := generateCharacterSprites(playerSheet, 256)
	playerSpawn := findOpenNode(levelData[0])
	playerActor := spawn_actor(playerSpawn.x, playerSpawn.y, "player", playerAnim)
	player = spawn_character(playerActor)
	player.maxhp = 40
	player.hp = 40
	player.actor.faction = friendly
	player.prange = 0.0
	player.arange = 5000.0

	characters = append(characters, player)

	terminalAnim := generateActorSprites(terminalSheet, 1, 128)
	terminalSpawn := findOpenNode(levelData[0])
	terminalActor := spawn_actor(terminalSpawn.x, terminalSpawn.y, "terminal", terminalAnim)
	terminalActor.direction = 3
	actors = append(actors, terminalActor)
	actors = append(actors, playerActor)
}

func lerp_64(v0x float64, v0y float64, v1x float64, v1y float64, t float64) (float64, float64) {
	return (1-t)*v0x + t*v1x, (1-t)*v0y + t*v1y
}

func inMapRange(x int, y int, levelData [2][32][32]*node) bool {
	if x >= 0 && x < len(levelData[0]) && y >= 0 && y < len(levelData[0]) {
		return true
	} else {
		return false
	}
}

func (g *Game) Update(screen *ebiten.Image) error {

	clearVisibility(levelData[0])
	compute_fov(vec{x: player.actor.x, y: player.actor.y}, levelData[0])

	for _, c := range characters {
		cx, cy := cartesianToIso(float64(c.actor.x), float64(c.actor.y))
		c.actor.coord.x, c.actor.coord.y = lerp_64(c.actor.coord.x, c.actor.coord.y, cx, cy, 0.06)
	}

	g.CamPosX, g.CamPosY = lerp_64(g.CamPosX, g.CamPosY, player.actor.coord.x, -player.actor.coord.y, 0.03)

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		// fmt.Printf("x: %d, y: %d", mx, my)
		player.actor.state = walk
		player.target = player

		tx, ty := getTileXY(g)
		if inMapRange(tx, ty, levelData) {

			for _, c := range characters {
				diffX, diffY := tx-player.target.actor.x, ty-player.target.actor.y
				if diffX*diffX < 2.0 && diffY*diffY < 4.0 && c != player {
					player.target = c
					break
				}
			}

			if tx < len(levelData[0]) && ty < len(levelData[0]) && tx >= 0 && ty >= 0 {
				player.dest = &node{x: tx, y: ty}
			}
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {

		tx, ty := getTileXY(g)
		if inMapRange(tx, ty, levelData) {
			//crashes if selection out of range
			if levelData[0][tx][ty].walkable {
				levelData[1][tx][ty].tile = rand.Intn(31)
			}
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.Key3) {

		tx, ty := getTileXY(g)
		if inMapRange(tx, ty, levelData) {
			creepAnim := generateCharacterSprites(creepSheet, 256)
			creepActor := spawn_actor(tx, ty, "creep", creepAnim)
			c := spawn_character(creepActor)
			c.dest = endOfTheRoad
			c.actor.faction = hostile
			c.prange = 8000.0
			c.arange = 5000.0
			actors = append(actors, creepActor)
			characters = append(characters, c)
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.Key4) {

		tx, ty := getTileXY(g)
		bossAnim := generateCharacterSprites(bossSheet, 512)
		if inMapRange(tx, ty, levelData) {
			bossActor := spawn_actor(tx, ty, "boss", bossAnim)
			c := spawn_character(bossActor)
			c.dest = endOfTheRoad
			c.actor.faction = hostile
			c.prange = 8000.0
			c.arange = 5000.0
			actors = append(actors, bossActor)
			characters = append(characters, c)
		}
	}

	if g.count%3 == 0 {
		characterStateMachine(characters, levelData[0])
		terminalStateMachine(actors)

		// terminalStateMachine(actors)
		actors, characters = removeDeadCharacters(actors, characters)

	}
	g.count++
	return nil
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

	return int(imx), int(imy)
}

type sprite struct {
	yi   float64
	pic  *ebiten.Image
	geom ebiten.GeoM
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x10, 0x10, 0x10, 1})
	// px, py := player.actor.coord.x, player.actor.coord.y

	var drawLater []*sprite

	for x := 0; x < len(levelData[0]); x++ {
		for y := 0; y < len(levelData[0]); y++ {
			if levelData[0][x][y].visible {
				xi, yi := cartesianToIso(float64(x), float64(y))
				// if player.actor.coord.y > yi {
				g.op.GeoM.Reset()
				g.op.GeoM.Translate(float64(xi), float64(yi))
				g.op.GeoM.Translate(-g.CamPosX, g.CamPosY)
				g.op.GeoM.Translate(float64(g.windowWidth/2.0), float64(g.windowHeight/2.0))

				t := tiles[levelData[0][x][y].tile]
				screen.DrawImage(t, g.op)

				d := doodads[levelData[1][x][y].tile]
				if levelData[1][x][y].tile > 0 {
					g.op.GeoM.Translate(-256.0, -400.0)
					screen.DrawImage(d, g.op)
					drawLater = append(drawLater, &sprite{yi: yi, pic: d, geom: g.op.GeoM})
				}
			}
		}
	}

	// DRAW ACTORS
	for _, a := range actors {

		isoSquare(g, screen, a.coord, 2.0, a.faction)

		startingFrame := 0
		// DRAW CHARACTER
		// the length of anims tells you if this is a character or item
		// characters will have an anims length of 6
		// widgets will have an anims length of 1
		startingFrame = a.direction * 10

		if levelData[0][a.x][a.y].visible {
			if len(a.anims) == 6 {
				g.op.GeoM.Reset()

				g.op.GeoM.Translate(float64(a.coord.x), float64(a.coord.y))
				if a.name == "boss" {
					g.op.GeoM.Translate(-224.0, -300.0)
				} else {
					g.op.GeoM.Translate(-96.0, -128.0)
				}
				g.op.GeoM.Translate(-g.CamPosX, g.CamPosY)
				g.op.GeoM.Translate(float64(g.windowWidth/2.0), float64(g.windowHeight/2.0))

				// The screen should be avoided as a render source
				// If I want the tiles to overlap the feet of the gopher, I'll need to
				// Create another render source for the gopher to prevent conflicting render calls
				// And then insert it into the painter's algorithm
				// screen.DrawImage(a.anims[a.state][(a.frame+startingFrame)], g.op)

				drawLater = append(drawLater, &sprite{yi: a.coord.y, pic: a.anims[a.state][(a.frame + startingFrame)], geom: g.op.GeoM})

				// targetRect(player.target.actor.coord, imd, player.target.actor.faction)
				// isoSquare(player.target.actor.coord, 3, imd, player.target.actor.faction)

			} else {
				g.op.GeoM.Reset()
				// DRAW WIDGETS
				g.op.GeoM.Translate(float64(a.coord.x), float64(a.coord.y))
				g.op.GeoM.Translate(-96.0, -96.0)
				g.op.GeoM.Translate(-g.CamPosX, g.CamPosY)
				g.op.GeoM.Translate(float64(g.windowWidth/2.0), float64(g.windowHeight/2.0))

				// raise y by 60 after i make a vec add fn
				//screen.DrawImage(a.anims[4][(a.frame+startingFrame)], g.op)

				drawLater = append(drawLater, &sprite{yi: a.coord.y, pic: a.anims[4][(a.frame + startingFrame)], geom: g.op.GeoM})

			}
		}
	}

	sort.Slice(drawLater[:], func(i, j int) bool {
		return drawLater[i].yi < drawLater[j].yi
	})

	for _, s := range drawLater {
		g.op.GeoM.Reset()
		g.op.GeoM = s.geom
		screen.DrawImage(s.pic, g.op)
	}

	// drawHealthPlates(g, screen, characters)

	// for _, a := range actors {

	// 	isoSquare(g, screen, a.coord, 2.0, a.faction)
	// }

	// Draw the sample text
	text.Draw(screen, sampleText, exocet_face, g.windowWidth-230, 30, color.White)

	ebitenutil.DebugPrint(
		screen,
		fmt.Sprintf("TPS: %0.2f\n", ebiten.CurrentTPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.windowWidth, g.windowHeight
}

func main() {
	g := &Game{
		Name:         "Diabgo",
		windowWidth:  1280,
		windowHeight: 900,
		tileSize:     64,
		CamPosX:      0,
		CamPosY:      0,
		CamSpeed:     500,
		op:           &ebiten.DrawImageOptions{},
	}

	ebiten.SetWindowSize(g.windowWidth, g.windowHeight)
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

func findOpenNode(levelData [32][32]*node) *node {
	x := rand.Intn(31)
	y := rand.Intn(31)
	for levelData[x][y].tile != 1 {
		x = rand.Intn(31)
		y = rand.Intn(31)
	}
	return &node{x: x, y: y}
}
