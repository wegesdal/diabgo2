package main

import "fmt"

type projectile struct {
	target    vec64
	actor     *actor
	max_range float64
	speed     float64
}

func spawnProjectile(start vec64, target vec64, actor *actor, max_range float64, speed float64) *projectile {
	return &projectile{target: target, actor: actor, max_range: max_range, speed: speed}
}

func projectileStateMachine(projectiles []*actor) {
	for _, p := range projectiles {
		fmt.Print(p.name)
		// lerp to destination
		//
	}
}
