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

const tileSize = 64.0

var (
	tilesImage   *ebiten.Image
	doodadsImage *ebiten.Image
	levelData    [2][mapSize][mapSize]*node
	// endOfTheRoad     *node
	tiles            []*ebiten.Image
	doodads          []*ebiten.Image
	healthGlobeImage *ebiten.Image
	manaGlobeImage   *ebiten.Image
	playerSheet      *ebiten.Image
	creepSheet       *ebiten.Image
	bossSheet        *ebiten.Image
	terminalSheet    *ebiten.Image
	healthGlobe      []*ebiten.Image
	manaGlobe        []*ebiten.Image
	player           *character
	actors           []*actor
	characters       []*character
	frames           int
	second           = time.Tick(time.Second)
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

type sprite struct {
	yi   float64
	pic  *ebiten.Image
	geom ebiten.GeoM
}

type vec64 struct {
	x float64
	y float64
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

	healthGlobeImage, _, err = ebitenutil.NewImageFromFile("assets/sprites/fs_sphere_red.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	manaGlobeImage, _, err = ebitenutil.NewImageFromFile("assets/sprites/fs_sphere_blue.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	//INITIALIZE WORLD OBJECTS

	// MAP
	levelData = generateMap()
	tiles = generateTiles(tilesImage)
	doodads = generateDoodads(doodadsImage)

	// ACTORS
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

	// UI
	healthGlobe = generateGlobeSprites(healthGlobeImage)
	manaGlobe = generateGlobeSprites(manaGlobeImage)

}

func lerp_64(v0x float64, v0y float64, v1x float64, v1y float64, t float64) (float64, float64) {
	return (1-t)*v0x + t*v1x, (1-t)*v0y + t*v1y
}

func inMapRange(x int, y int, levelData [2][mapSize][mapSize]*node) bool {
	if x >= 0 && x < len(levelData[0]) && y >= 0 && y < len(levelData[0]) {
		return true
	} else {
		return false
	}
}

func (g *Game) Update(screen *ebiten.Image) error {

	if g.count%2 == 0 {
		vision_range := 6.0
		head_room := int(math.Min(float64(player.actor.x), vision_range))
		foot_room := int(math.Min(mapSize-float64(player.actor.x)-1, vision_range))
		width := head_room + foot_room
		grid_to_check := make([][]*node, width)
		minY := int(math.Max(0, float64(player.actor.y)-vision_range))
		maxY := int(math.Min(mapSize-1, float64(player.actor.y)+vision_range))
		for i := 0; i < width; i++ {
			grid_to_check[i] = levelData[0][player.actor.x+i-head_room][minY:maxY]
		}
		clearVisibility(levelData[0])

		compute_fov(vec{x: head_room, y: int(math.Min(float64(player.actor.y), vision_range))}, grid_to_check)
	}
	for _, c := range characters {
		cx, cy := cartesianToIso(float64(c.actor.x), float64(c.actor.y))
		c.actor.coord.x, c.actor.coord.y = lerp_64(c.actor.coord.x, c.actor.coord.y, cx, cy, 0.06)
	}

	g.CamPosX, g.CamPosY = lerp_64(g.CamPosX, g.CamPosY, player.actor.coord.x, -player.actor.coord.y, 0.03)

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		// fmt.Printf("x: %d, y: %d", mx, my)
		player.actor.state = walk
		player.target = player

		//Get the cursor position
		mx, my := ebiten.CursorPosition()
		//Offset for center
		fmx := float64(mx) - float64(g.windowWidth)/2.0
		fmy := float64(my) - float64(g.windowHeight)/2.0

		x, y := fmx+g.CamPosX, fmy-g.CamPosY

		offset := vec64{x: 32.0, y: -30.0}

		for _, c := range characters {
			diff := sub_vec64(c.actor.coord, vec64{x: x - offset.x, y: y - offset.y})
			if math.Abs(diff.x) < 64 && math.Abs(diff.y) < 64 && c != player {
				player.target = c
				break
			}

		}
		tx, ty := getTileXY(g)
		if inMapRange(tx, ty, levelData) {

			// prevent movement to invisible tiles
			// for performance reasons (pathfinding chokes)
			if levelData[0][tx][ty].visible && levelData[0][tx][ty].walkable {
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
			// c.dest = endOfTheRoad
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
			// c.dest = endOfTheRoad
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

	for _, c := range characters {
		if c == player.target {
			isoSquare(g, screen, c.actor.coord, 2.0, c.actor.faction)
			isoTargetDebug(g, screen, c.actor.coord)
		}
	}

	// DRAW ACTORS
	for _, a := range actors {
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

	// UI
	// Draw the sample text
	text.Draw(screen, sampleText, exocet_face, g.windowWidth-230, 30, color.White)

	// TODO: LERP THE RESOURCE GLOBE FRAMES
	if player.hp > 0 {
		g.op.GeoM.Reset()
		percentHealth := player.hp * 100 / player.maxhp
		pixelsPerFrame := int(128 / 25)
		healthGlobeFrame := 25 - (int(percentHealth / 4))
		g.op.GeoM.Translate(0.0, float64(g.windowHeight+healthGlobeFrame*pixelsPerFrame)-128)
		screen.DrawImage(healthGlobe[healthGlobeFrame], g.op)

		g.op.GeoM.Reset()
		percentMana := player.hp * 100 / player.maxhp
		manaGlobeFrame := 25 - (int(percentMana / 4))
		g.op.GeoM.Translate(float64(g.windowWidth)-128.0, float64(g.windowHeight+healthGlobeFrame*pixelsPerFrame)-128)
		screen.DrawImage(manaGlobe[manaGlobeFrame], g.op)

		ebitenutil.DebugPrint(
			screen,
			fmt.Sprintf("TPS: %0.2f\n", ebiten.CurrentTPS()))

	}
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
