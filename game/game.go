package game

import (
	"fmt"
	"image/color"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"

	"github.com/willybeans/bumble_bash/assets"
)

const (
	screenWidth  = 800
	screenHeight = 600

	flowerSpawnTime  = 1 * time.Second
	hoseSpawnTime    = 1 * time.Second
	dropletSpawnTime = 2 * time.Second

	baseFlowerVelocity  = 0.25
	flowerSpeedUpAmount = 0.1
	flowerSpeedUpTime   = 5 * time.Second
)

var isHit = false

type Game struct {
	mut               sync.Mutex
	player            *Player
	flowerSpawnTimer  *Timer
	hoseSpawnTimer    *Timer
	dropletSpawnTimer *Timer

	flowers  []*Flower
	hoses    []*Hose
	droplets []*Droplet

	score int

	baseVelocity  float64
	velocityTimer *Timer
}

func NewGame() *Game {
	g := &Game{
		dropletSpawnTimer: NewTimer(dropletSpawnTime),
		flowerSpawnTimer:  NewTimer(flowerSpawnTime),
		hoseSpawnTimer:    NewTimer(hoseSpawnTime),
		baseVelocity:      baseFlowerVelocity,
		velocityTimer:     NewTimer(flowerSpeedUpTime),
	}

	g.player = NewPlayer(g)

	return g
}

func (g *Game) Update() error {
	g.velocityTimer.Update()
	if g.velocityTimer.IsReady() {
		g.velocityTimer.Reset()
		g.baseVelocity += flowerSpeedUpAmount
	}

	g.player.Update()

	g.flowerSpawnTimer.Update()
	if g.flowerSpawnTimer.IsReady() {
		g.flowerSpawnTimer.Reset()

		f := NewFlower(g.baseVelocity)
		g.flowers = append(g.flowers, f)
	}

	g.hoseSpawnTimer.Update()
	if g.hoseSpawnTimer.IsReady() {
		g.hoseSpawnTimer.Reset()

		if len(g.hoses) < 2 {
			h := NewHose(g, len(g.hoses))
			g.hoses = append(g.hoses, h)
		}
	}

	for _, f := range g.flowers {
		f.Update()
	}

	for _, d := range g.droplets {
		d.Update()
	}

	for _, h := range g.hoses {
		h.Update()
	}

	for i, f := range g.flowers {
		isPlayerCollision := f.Collider().Intersects(g.player.Collider())

		if f.position.X < 0 ||
			f.position.Y < 0 ||
			f.position.X > screenWidth ||
			f.position.Y > screenHeight ||
			isPlayerCollision {

			g.flowers = RemoveIndex(g.flowers, i, f)

			if isPlayerCollision {
				g.score++
				p := NewPollen(g.baseVelocity, g.player.sprite)
				g.player.Pollens = append(g.player.Pollens, p)
			}
		}

	}

	for i, d := range g.droplets {
		isPlayerCollision := d.Collider().Intersects(g.player.Collider())
		if d.position.X < 0 ||
			d.position.Y < 0 ||
			d.position.X > screenWidth ||
			d.position.Y > screenHeight ||
			isPlayerCollision {

			g.droplets = RemoveIndex(g.droplets, i, d)

			if isPlayerCollision {
				g.player.isHit = true
				g.player.hitCoords = Vector{
					X: d.Collider().X,
					Y: d.Collider().Y,
				}
				g.score--

				if len(g.player.Pollens) > 0 {
					//needs to be connected with the pollen struct for falling effect
					// var random = randIntRange(0, len(g.player.Pollens))
					// g.player.Pollens = g.player.Pollens[:len(g.player.Pollens)-random]
				}
			}

		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.player.Draw(screen)

	for _, f := range g.flowers {
		f.Draw(screen)
	}
	for _, h := range g.hoses {
		h.Draw(screen)
	}
	for _, d := range g.droplets {
		d.Draw(screen)
	}
	for _, p := range g.player.Pollens {
		p.Draw(screen, g.player.position)
	}

	text.Draw(screen, fmt.Sprintf("%06d", g.score), assets.ScoreFont, screenWidth/2-100, 50, color.White)

	// for _, f := range g.flowers {
	// 	vector.StrokeRect(
	// 		screen,
	// 		float32(f.position.X),
	// 		float32(f.position.Y),
	// 		float32(f.sprite.Bounds().Dx()),
	// 		float32(f.sprite.Bounds().Dy()),
	// 		1.0,
	// 		color.White,
	// 		false,
	// 	)
	// }

	// for _, d := range g.droplets {
	// 	vector.StrokeRect(
	// 		screen,
	// 		float32(d.position.X),
	// 		float32(d.position.Y),
	// 		float32(d.sprite.Bounds().Dx()),
	// 		float32(d.sprite.Bounds().Dy()),
	// 		1.0,
	// 		color.White,
	// 		false,
	// 	)
	// }

	// vector.StrokeRect(
	// 	screen,
	// 	float32(g.player.position.X),
	// 	float32(g.player.position.Y),
	// 	float32(g.player.sprite.Bounds().Dx()),
	// 	float32(g.player.sprite.Bounds().Dy()),
	// 	1.0,
	// 	color.White,
	// 	false,
	// )

	// for _, h := range g.hoses {
	// 	vector.StrokeRect(
	// 		screen,
	// 		float32(h.position.X),
	// 		float32(h.position.Y),
	// 		float32(h.sprite.Bounds().Dx()),
	// 		float32(h.sprite.Bounds().Dy()),
	// 		1.0,
	// 		color.White,
	// 		false,
	// 	)
	// }

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) AddDroplet(d *Droplet) {
	g.droplets = append(g.droplets, d)
}
