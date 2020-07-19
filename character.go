package main

import (
	"math"
)

type character struct {
	actor  *actor
	maxhp  int
	dest   *node
	hp     int
	prange float64
	arange float64
	target *character
}

const (
	light = iota
	dark
)

func spawn_character(a *actor) *character {

	var c = character{actor: a}

	c.dest = &node{x: a.x, y: a.y}
	c.maxhp = 20
	c.hp = 10
	c.target = &c
	return &c
}

func step_forward(a *actor, path []*node) {
	if len(path) > 0 {
		ix, iy := isoToCartesian(a.coord.x, a.coord.y)
		// don't update next block until close
		if math.Pow(ix-float64(a.x), 2.0)+math.Pow(iy-float64(a.y), 2.0) < 0.5 {

			a.direction = wayfind(a.x, a.y, path[len(path)-1].x, path[len(path)-1].y)
			a.x = path[len(path)-1].x
			a.y = path[len(path)-1].y
		}
	}
}

func characterStateMachine(characters []*character, levelData [mapSize][mapSize]*node) {

	for _, c := range characters {

		// advance the animation frame
		c.actor.frame = (c.actor.frame + 1) % 10

		// auto targeting for non-players
		for _, o := range characters {

			// _, ocharmed := o.actor.effects["charmed"]

			// oc := 1
			// if ocharmed {
			// 	oc = -1
			// 	o.hp += 1
			// }

			// friendly is positive, enemy is negative, neutral is 0
			// if both are friendly the product of their states is positive
			// if both are hostile the product of their states is positive
			// if one is neutral the product of their states is 0
			// if they are opposed the product of their states is negative

			if c != o && c.actor.state != dead && o.actor.state != dead {
				// let the player target manually
				if c.actor.name != "player" {
					d := sub_vec64(c.actor.coord, o.actor.coord)
					d_square := d.x*d.x + d.y*d.y
					if (d_square < c.prange || d_square < c.arange) && o.actor.faction*c.actor.faction < 0 {
						c.target = o
						break
					}
				}
			}
		}
	}

	for _, c := range characters {
		d := sub_vec64(c.target.actor.coord, c.actor.coord)
		d_square := d.x*d.x + d.y*d.y

		ix, iy := isoToCartesian(c.actor.coord.x, c.actor.coord.y)
		if c.actor.state == idle {

			// if actor has not reached destination, walk

			if c.dest.x != int(ix+0.5) && c.dest.y != int(iy+0.5) || c.target != c {
				c.actor.state = walk
			}

		} else if c.actor.state == walk {

			// if actor has reached destination, idle

			if c.dest.x == int(ix+0.5) && c.dest.y == int(iy+0.5) {
				c.actor.state = idle
			}

			if c.target != c {
				// if in range, attack
				if d_square < c.arange {
					c.actor.state = attack
					c.actor.direction = wayfind(c.actor.x, c.actor.y, c.target.actor.x, c.target.actor.y)
				} else {
					// otherwise move towards target unless player (let the player control their movement)

					// I SHOULDN'T RECALCULATE EVERY LOOP
					path := Astar(&node{x: c.actor.x, y: c.actor.y}, &node{x: c.target.actor.x, y: c.target.actor.y}, levelData, true)
					if len(path) > 0 {
						if path[len(path)-1].x != c.target.actor.x || path[len(path)-1].y != c.target.actor.y {
							step_forward(c.actor, path)
						}
					}
				}
				// if no target
			} else {
				path := Astar(&node{x: c.actor.x, y: c.actor.y}, c.dest, levelData, true)
				step_forward(c.actor, path)
			}

		} else if c.actor.state == attack {
			if d_square < c.arange {
				if c.actor.frame == 9 {
					c.target.hp -= 3
				}
			} else {
				c.actor.state = idle
			}

			if c.target.hp < 1 {
				c.actor.state = idle
				c.target.actor.frame = 0
				c.target.actor.state = dead
				c.target = c
			}
		}
	}
}

// func drawHealthPlates(g *Game, screen *ebiten.Image, characters []*character) {
// 	for _, c := range characters {
// 		// total length of health plate
// 		length := 40.0
// 		// number of bars to represent health (10 hp per bar)
// 		bars := c.maxhp / 10
// 		// length of a single bar
// 		bar_length := length / float64(bars)
// 		start_X := c.actor.coord.x - g.CamPosX + 1280/2

// 		if c.hp > 0 {
// 			for i := 0; i < bars; i++ {
// 				verticalOffset := g.CamPosY + 300
// 				x1 := start_X + float64(i)*bar_length + 1
// 				y := c.actor.coord.y + verticalOffset
// 				if i*10 <= c.hp && (i+1)*10 > c.hp {
// 					f := float64(10-c.hp%10) / 10
// 					x2 := start_X + float64(i+1)*bar_length - f*bar_length
// 					ebitenutil.DrawLine(screen, x1, y, x2, y, factionColor(c.actor.faction, light))
// 					x1 = x2
// 					x2 = start_X + float64(i+1)*bar_length - 1
// 					ebitenutil.DrawLine(screen, x1, y, x2, y, factionColor(c.actor.faction, dark))
// 				} else {
// 					// draw the whole bar
// 					x2 := start_X + float64(i+1)*bar_length - 1
// 					ebitenutil.DrawLine(screen, x1, y, x2, y, factionColor(c.actor.faction, light))
// 				}
// 			}
// 		}

// 	}

// }

func removeDeadActors(c *character, actors []*actor) []*actor {
	for j, a := range actors {
		// remove the actor from the actors slice first
		if a == c.actor {
			actors[j] = actors[len(actors)-1]
			actors[len(actors)-1] = nil
			actors = actors[:len(actors)-1]
		}
	}
	return actors
}

func removeDeadCharacters(actors []*actor, characters []*character) ([]*actor, []*character) {

	for i, c := range characters {
		// KILL CREEPS WHO REACH END OF THE ROAD
		// TODO: ADJUST SCORE
		// if c.actor.name == "creep" && c.actor.x == c.dest.x && c.actor.y == c.dest.y && c.actor.state != dead {
		// 	c.actor.frame = 0
		// 	c.actor.state = dead
		// }
		if c.actor.state == dead && c.actor.frame == 9 {

			actors = removeDeadActors(c, actors)

			// remove the character from the character slice
			characters[i] = characters[len(characters)-1]
			characters[len(characters)-1] = nil
			characters = characters[:len(characters)-1]
			// break out of slice (dangerous to continue to modify a slice while iterating it)
			break
		}
	}
	return actors, characters
}
