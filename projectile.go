package main

import (
	"container/list"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
)

type projectile struct {
	coord     vec64
	target    vec64
	particles *list.List
	max_range float64
	speed     float64
}

func spawnProjectile(coord vec64, target vec64, max_range float64, speed float64) *projectile {
	p := &projectile{coord: coord, target: target, max_range: max_range, speed: speed}
	p.particles = list.New()
	return p
}

func newParticle(img *ebiten.Image) *particle {
	c := rand.Intn(50) + 300
	dir := rand.Float64() * 2 * math.Pi
	a := rand.Float64() * 2 * math.Pi
	s := rand.Float64()*0.1 + 0.4
	return &particle{
		img:      img,
		maxCount: c,
		count:    c,
		dir:      dir,
		angle:    a,
		scale:    s,
		alpha:    0.5,
	}
}

func (p *particle) update() {
	if p.count == 0 {
		return
	}
	p.count--
}

func (p *particle) terminated() bool {
	return p.count == 0
}

func projectileStateMachine(projectiles []*projectile) {
	for _, p := range projectiles {
		if p.max_range > 0 {

			if p.particles.Len() < 100 && rand.Intn(4) < 3 {
				// EMIT
				p.particles.PushBack(newParticle(smokeImage))
			}
			p.max_range--
		}

		for e := p.particles.Front(); e != nil; e = e.Next() {
			s := e.Value.(*particle)
			s.update()
			if s.terminated() {
				defer p.particles.Remove(e)
			}
		}

	}
}

type particle struct {
	count    int
	maxCount int
	dir      float64
	img      *ebiten.Image
	op       *ebiten.DrawImageOptions
	scale    float64
	angle    float64
	alpha    float64
}

func drawProjectiles(g *Game, screen *ebiten.Image) {
	for _, p := range projectiles {
		for e := p.particles.Front(); e != nil; e = e.Next() {
			particle := e.Value.(*particle)
			x := math.Cos(particle.dir) * float64(particle.maxCount-particle.count)
			y := math.Sin(particle.dir) * float64(particle.maxCount-particle.count)
			op := &ebiten.DrawImageOptions{}
			sx, sy := particle.img.Size()
			op.GeoM.Translate(-float64(sx)/2, -float64(sy)/2)
			op.GeoM.Rotate(particle.angle)
			op.GeoM.Scale(particle.scale, particle.scale)
			op.GeoM.Translate(x, y)
			op.GeoM.Translate(p.coord.x, p.coord.y)

			op.GeoM.Translate(-g.CamPosX, g.CamPosY)
			op.GeoM.Translate(float64(g.windowWidth/2.0), float64(g.windowHeight/2.0))

			rate := float64(particle.count) / float64(particle.maxCount)
			alpha := 0.0
			if rate < 0.2 {
				alpha = rate / 0.2
			} else if rate > 0.8 {
				alpha = (1 - rate) / 0.2
			} else {
				alpha = 1
			}
			alpha *= particle.alpha
			op.ColorM.Scale(1, 1, 1, alpha)
			particle.op = op
			screen.DrawImage(particle.img, particle.op)
		}
	}
}
