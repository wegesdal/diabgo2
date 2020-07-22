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

const tileSize = 128.0

var (
	gx               int
	gy               int
	tilesImage       *ebiten.Image
	doodadsImage     *ebiten.Image
	levelData        [3][3][chunkSize][chunkSize]*node
	flatMap          [3 * chunkSize][3 * chunkSize]*node
	gradient         [gradSize][gradSize][2]float64
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
	projectiles      []*projectile
	frames           int
	second           = time.Tick(time.Second)
	bossAnim         map[int][]*ebiten.Image
	creepAnim        map[int][]*ebiten.Image
	smokeImage       *ebiten.Image
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
	tilesImage, _, err = ebitenutil.NewImageFromFile("assets/sprites/blocks.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	playerSheet, _, err = ebitenutil.NewImageFromFile("assets/sprites/gopher.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	creepSheet, _, err = ebitenutil.NewImageFromFile("assets/sprites/fs_creep.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	bossSheet, _, err = ebitenutil.NewImageFromFile("assets/sprites/diabgopher.png", ebiten.FilterDefault)
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

	smokeImage, _, err = ebitenutil.NewImageFromFile("assets/sprites/smoke.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	//INITIALIZE WORLD OBJECTS

	// MAP
	levelData = generateMap()
	flatMap = flattenMap()
	tiles = generateTiles(tilesImage)
	doodads = generateDoodads(doodadsImage)

	gradient = generateGradient()

	// ACTORS
	playerAnim := generateCharacterSprites(playerSheet, 256)
	playerSpawn := findOpenNode()
	playerActor := spawn_actor(playerSpawn.x, playerSpawn.y, "player", playerAnim)
	player = spawn_character(playerActor)
	player.maxhp = 40
	player.hp = 40
	player.actor.faction = friendly
	player.prange = 0.0
	player.arange = 5000.0

	characters = append(characters, player)

	terminalAnim := generateActorSprites(terminalSheet, 1, 128)
	terminalSpawn := findOpenNode()
	terminalActor := spawn_actor(terminalSpawn.x, terminalSpawn.y, "terminal", terminalAnim)
	terminalActor.direction = 3
	actors = append(actors, terminalActor)
	actors = append(actors, playerActor)

	// UI
	healthGlobe = generateGlobeSprites(healthGlobeImage)
	manaGlobe = generateGlobeSprites(manaGlobeImage)

	creepAnim = generateCharacterSprites(creepSheet, 256)
	bossAnim = generateCharacterSprites(bossSheet, 512)

	gx = 1
	gy = 1

}

func lerp_64(v0x float64, v0y float64, v1x float64, v1y float64, t float64) (float64, float64) {
	return (1-t)*v0x + t*v1x, (1-t)*v0y + t*v1y
}

func globalToMapCoords(x int, y int) (int, int) {
	return int(x) + (1-gx)*chunkSize, int(y) + (1-gy)*chunkSize
}

func (g *Game) Update(screen *ebiten.Image) error {
	if g.count%5 == 0 {

		// INTERPOLATE PROJECTILE MOVEMENT
		for _, p := range projectiles {
			p.coord.x, p.coord.y = lerp_64(p.coord.x, p.coord.y, p.target.x, p.target.y, p.speed)
			if p.timer > 0 {
				if p.name == "smoke" {
					for _, c := range characters {
						if math.Abs(float64(c.actor.coord.x-p.coord.x))+math.Abs(float64(c.actor.coord.y-p.coord.y)) < 50.0 {
							c.hp--

							if c.hp == 0 {
								c.actor.frame = 0
								c.actor.state = dead
							}
						}
					}
					// ix, iy := isoToCartesian(p.coord.x, p.coord.y)
					// proj_map_x, proj_map_y := globalToMapCoords(int(ix), int(iy))
					// if proj_map_x > 0 && proj_map_y > 0 && proj_map_x < chunkSize*3-1 && proj_map_y < chunkSize*3-1 {
					// 	// flatMap[gx][gy].blocks_vision = true
					// }
				}
			}
		}

		for chunk_x := 0; chunk_x < 3; chunk_x++ {
			for chunk_y := 0; chunk_y < 3; chunk_y++ {
				clearVisibility(levelData[chunk_x][chunk_y])
			}
		}
		player_map_x, player_map_y := globalToMapCoords(player.actor.x, player.actor.y)
		compute_fov(vec{x: player_map_x, y: player_map_y}, flatMap)
	}

	for _, c := range characters {
		// INTERPOLATE CHARACTER MOVEMENT
		cx, cy := cartesianToIso(float64(c.actor.x), float64(c.actor.y))
		c.actor.coord.x, c.actor.coord.y = lerp_64(c.actor.coord.x, c.actor.coord.y, cx, cy, 0.06)
	}

	g.CamPosX, g.CamPosY = lerp_64(g.CamPosX, g.CamPosY, player.actor.coord.x, -player.actor.coord.y, 0.03)

	// LEFT MOUSE (PRIMARY ACTION)
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {

		player.actor.state = walk
		player.target = player

		//Get the cursor position
		mx, my := ebiten.CursorPosition()
		//Offset for center
		fmx := float64(mx) - float64(g.windowWidth)/2.0
		fmy := float64(my) - float64(g.windowHeight)/2.0

		// TARGETING BOX TRIGGER
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
		player.dest = &node{x: tx, y: ty}
	}

	if inpututil.IsKeyJustPressed(ebiten.Key2) {

		// tx, ty := getTileXY(g)
		// if inMapRange(tx, ty, levelData[1]) {
		// 	//crashes if selection out of range
		// 	if levelData[1][1][tx][ty].walkable {
		// 		levelData[1][1][tx][ty].tile = rand.Intn(31)
		// 	}
		// }
	}

	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		tx, ty := getTileXY(g)
		spawnCreep(tx, ty)
	}

	if inpututil.IsKeyJustPressed(ebiten.Key4) {
		tx, ty := getTileXY(g)

		spawnBoss(tx, ty)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		player.hp_target += 10
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		player.hp_target += 10
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		//Get the cursor position

		// POINT SPELL CAST IN CORRECT DIRECTION
		player.actor.state = cast
		player.actor.frame = 0
		mx, my := ebiten.CursorPosition()
		fmx := float64(mx) - float64(g.windowWidth)/2.0
		fmy := float64(my) - float64(g.windowHeight)/2.0
		x, y := fmx+g.CamPosX, fmy-g.CamPosY
		radians := math.Atan2(float64(fmx), float64(fmy))
		switch degrees := 180 + radians*180/3.14159265; {
		case degrees > 315:
			player.actor.direction = 5
		case degrees > 270:
			player.actor.direction = 6
		case degrees > 225:
			player.actor.direction = 7
		case degrees > 180:
			player.actor.direction = 0
		case degrees > 135:
			player.actor.direction = 1
		case degrees > 90:
			player.actor.direction = 2
		case degrees > 45:
			player.actor.direction = 3
		case degrees > 0:
			player.actor.direction = 4
		}
		projectiles = append(projectiles, spawnProjectile(player.actor.coord, vec64{x: x, y: y}, 30.0, 0.2, "smoke"))
	}

	if g.count%3 == 0 {
		characterStateMachine(characters)
		terminalStateMachine(actors)
		projectileStateMachine(projectiles)

	}
	g.count++

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x10, 0x10, 0x10, 1})

	chunk()

	var drawLater []*sprite
	for cx := 0; cx < 3; cx++ {
		for cy := 0; cy < 3; cy++ {
			for x := 0; x < len(levelData[cx][cy]); x++ {
				for y := 0; y < len(levelData[cx][cy]); y++ {
					if levelData[cx][cy][x][y].visible {
						if levelData[cx][cy][x][y].tile == sentinal {

							levelData[cx][cy][x][y].tile, levelData[cx][cy][x][y].walkable, levelData[cx][cy][x][y].blocks_vision = compute_noise(levelData[cx][cy][x][y].x, levelData[cx][cy][x][y].y)
							levelData[cx][cy][x][y].visible = true
						} else {
							xi, yi := cartesianToIso(float64(levelData[cx][cy][x][y].x), float64(levelData[cx][cy][x][y].y))

							g.op.GeoM.Reset()
							g.op.GeoM.Translate(float64(xi), float64(yi))
							g.op.GeoM.Translate(-g.CamPosX, g.CamPosY)
							g.op.GeoM.Translate(float64(g.windowWidth/2.0), float64(g.windowHeight/2.0))

							t := tiles[levelData[cx][cy][x][y].tile-1]
							// t := tiles[cx+cy*3]
							screen.DrawImage(t, g.op)

							// d := doodads[levelData[cx][cy][1][x][y].tile]
							// if levelData[cx][cy][1][x][y].tile > 0 {
							// 	g.op.GeoM.Translate(-256.0, -400.0)
							// 	screen.DrawImage(d, g.op)
							// 	drawLater = append(drawLater, &sprite{yi: yi, pic: d, geom: g.op.GeoM})
							// }
						}
					}
				}
			}
		}
	}

	for _, c := range characters {
		if c == player.target {
			isoSquare(g, screen, c.actor.coord, c.actor.faction)
			isoTargetDebug(g, screen, c.actor.coord)
		}
	}

	// DRAW ACTORS
	for _, a := range actors {
		startingFrame := 0
		// DRAW CHARACTER
		// the length of anims tells you if this is a character or widget
		// characters will have an anims length of 6
		// widgets will have an anims length of 1
		startingFrame = a.direction * 10

		x, y := globalToMapCoords(a.x, a.y)

		if x > 0 && x < chunkSize*3-1 && y > 0 && y < chunkSize*3-1 {
			if flatMap[x][y].visible {
				if len(a.anims) == 6 {
					g.op.GeoM.Reset()

					g.op.GeoM.Translate(float64(a.coord.x), float64(a.coord.y))
					if a.name == "boss" {
						g.op.GeoM.Translate(-224.0, -300.0)
					} else {
						g.op.GeoM.Translate(-72.0, -96.0)
					}
					g.op.GeoM.Translate(-g.CamPosX, g.CamPosY)
					g.op.GeoM.Translate(float64(g.windowWidth/2.0), float64(g.windowHeight/2.0))

					// The screen should be avoided as a render source
					// If I want the tiles to overlap the feet of the gopher, I'll need to
					// Create another render source for the gopher to prevent conflicting render calls
					// And then insert it into the painter's algorithm
					// screen.DrawImage(a.anims[a.state][(a.frame+startingFrame)], g.op)

					drawLater = append(drawLater, &sprite{yi: a.coord.y, pic: a.anims[a.state][(a.frame + startingFrame)], geom: g.op.GeoM})

				} else {
					g.op.GeoM.Reset()

					// DRAW WIDGETS
					g.op.GeoM.Translate(float64(a.coord.x), float64(a.coord.y))
					g.op.GeoM.Translate(-96.0, -96.0)
					g.op.GeoM.Translate(-g.CamPosX, g.CamPosY)
					g.op.GeoM.Translate(float64(g.windowWidth/2.0), float64(g.windowHeight/2.0))

					drawLater = append(drawLater, &sprite{yi: a.coord.y, pic: a.anims[4][(a.frame + startingFrame)], geom: g.op.GeoM})

				}
			} else {
				a.state = idle
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

	// DRAW PROJECTILE PARTICLES
	drawProjectiles(g, screen)

	// drawHealthPlates(g, screen, characters)

	// UI
	// Draw the sample text
	text.Draw(screen, sampleText, exocet_face, g.windowWidth-230, 30, color.White)

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
			fmt.Sprintf("TPS: %0.2f\ngx: %d\ngy: %d\n", ebiten.CurrentTPS(), gx, gy))

	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.windowWidth, g.windowHeight
}

func main() {
	g := &Game{
		Name:         "Diabgo",
		windowWidth:  1280,
		windowHeight: 760,
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
