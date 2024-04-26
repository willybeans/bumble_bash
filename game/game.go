package game

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"

	"github.com/willybeans/bumble_bash/assets"
)

const (
	screenWidth  = 800
	screenHeight = 600

	flowerSpawnTime = 1 * time.Second

	baseFlowerVelocity  = 0.25
	flowerSpeedUpAmount = 0.1
	flowerSpeedUpTime   = 5 * time.Second
)

var fallSpeed = 1.0

type Game struct {
	player           *Player
	flowerSpawnTimer *Timer
	flowers          []*Flower
	droplets         []*Droplet

	score int

	baseVelocity  float64
	velocityTimer *Timer
}

func NewGame() *Game {
	g := &Game{
		flowerSpawnTimer: NewTimer(flowerSpawnTime),
		baseVelocity:     baseFlowerVelocity,
		velocityTimer:    NewTimer(flowerSpeedUpTime),
	}

	g.player = NewPlayer(g)
	// g.flowers = [NewFlower(0.25)]

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

		d := NewDroplet(g.baseVelocity)
		g.droplets = append(g.droplets, d)

	}

	for _, f := range g.flowers {
		f.Update()
	}

	for _, d := range g.droplets {
		d.Update()
	}

	for i, f := range g.flowers {
		if f.Collider().Intersects(g.player.Collider()) {
			g.flowers = append(g.flowers[:i], g.flowers[i+1:]...)
			g.score++
		}
	}

	for i, d := range g.droplets {
		if d.Collider().Intersects(g.player.Collider()) {
			g.droplets = append(g.droplets[:i], g.droplets[i+1:]...)
			//make random number
			g.score--
		}
	}

	// if g.flower.Collider().Intersects(g.player.Collider()) {
	// 	println("player: ", g.player.position.X < g.flower.position.X)
	// }

	// if g.flower.Collider().Intersects(g.player.Collider()) {
	// 	g.flower.Update()
	// }
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, f := range g.flowers {
		f.Draw(screen)
	}
	for _, d := range g.droplets {
		d.Draw(screen)
	}
	// g.flower.Draw(screen)
	g.player.Draw(screen)

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

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
