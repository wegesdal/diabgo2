package main

func terminalStateMachine(actors []*actor) {

	var player *actor
	for _, p := range actors {
		if p.name == "player" {
			player = p
		}
	}

	activation_radius := 2000.0

	for _, a := range actors {
		if a.name == "terminal" {

			if player != nil {
				d := sub_vec64(player.coord, a.coord)
				d_square := d.x*d.x + d.y*d.y

				// fmt.Println(d_square)

				if d_square < activation_radius {
					if player.state == idle {
						a.state = activate
						player.state = activate
						player.x = a.x
						player.y = a.y
						// jostling the coord below makes the depth sort put the player in front
						player.coord.y += 15.0
						player.direction = 4
						a.frame = 0
					} else {
						if a.frame < 9 {
							a.frame++
						}
					}
				} else {
					if a.frame > 0 {
						a.frame--
					}
				}
			}
		}
	}
}

// func handleTerminalInput(player *character, txt *text.Text, input string) string {
// 	if player.actor.state == activate {
// 		txt.WriteString(win.Typed())
// 		input += win.Typed()
// 		if win.JustPressed(pixelgl.KeyEnter) || win.Repeated(pixelgl.KeyEnter) {
// 			switch input {
// 			case "foo":
// 				txt.WriteRune('\n')
// 				txt.WriteString("bar")
// 			case "heal":
// 				player.hp = player.maxhp
// 				txt.WriteRune('\n')
// 				txt.WriteString("completed")
// 			default:
// 				txt.WriteRune('\n')
// 				txt.WriteString(input + " is not defined")
// 			}
// 			input = ""
// 			txt.WriteRune('\n')
// 			txt.WriteString("> ")

// 		}
// 	}
// 	return input
// }

// func renderTerminalText(player *character, txt *text.Text, input string) string {

// 	if player.actor.state == activate {
// 		txt.Draw(win, pixel.IM.Moved(pixel.Vec{X: player.actor.coord.X + 23.0, Y: player.actor.coord.Y + 90.0}.Sub(txt.Bounds().Min)))
// 	} else {
// 		input = ""
// 		txt.Clear()
// 		txt.Color = colornames.Yellow
// 		txt.WriteString("> ")
// 	}
// 	return input
// }
