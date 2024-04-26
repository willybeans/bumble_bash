package game

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/willybeans/bumble_bash/assets"
)

// const (
// 	rotationSpeedMin = -0.02
// 	rotationSpeedMax = 0.02
// )

type Droplet struct {
	position      Vector
	rotation      float64
	movement      Vector
	rotationSpeed float64
	sprite        *ebiten.Image
}

func NewDroplet(baseVelocity float64) *Droplet {
	// sprite := assets.DropletSprite
	// bounds := sprite.Bounds()
	// halfW := float64(bounds.Dx()) / 2
	// halfH := float64(bounds.Dy()) / 2

	// target := Vector{
	// 	X: screenWidth * .75,
	// 	Y: screenHeight * .75,
	// }

	target := Vector{
		X: screenWidth / 2,
		Y: screenHeight / 2,
	}

	angle := rand.Float64() * 2 * math.Pi
	radius := screenWidth / 2.0

	pos := Vector{
		X: target.X + math.Cos(angle)*radius,
		Y: target.Y + math.Sin(angle)*radius,
	}
	// pos := Vector{
	// 	X: target.X + 100,
	// 	Y: target.Y,
	// }

	velocity := baseVelocity + rand.Float64()*1.5

	direction := Vector{
		X: target.X - pos.X,
		Y: target.Y - pos.Y,
	}
	normalizedDirection := direction.Normalize()

	movement := Vector{
		X: normalizedDirection.X * velocity,
		Y: normalizedDirection.Y * velocity,
	}

	sprite := assets.DropletSprite

	r := &Droplet{
		position:      pos,
		movement:      movement,
		rotationSpeed: rotationSpeedMin + rand.Float64()*(rotationSpeedMax-rotationSpeedMin),
		sprite:        sprite,
	}
	return r
}

func (d *Droplet) Update() {
	d.position.X += d.movement.X
	d.position.Y += d.movement.Y
	d.rotation += d.rotationSpeed
}

func (d *Droplet) Draw(screen *ebiten.Image) {
	// op := &ebiten.DrawImageOptions{}
	// // op.GeoM.Rotate(d.rotation)
	// op.GeoM.Translate(d.position.X, d.position.Y)
	// screen.DrawImage(d.sprite, op)

	bounds := d.sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-halfW, -halfH)
	op.GeoM.Rotate(d.rotation)
	op.GeoM.Translate(halfW, halfH)

	op.GeoM.Translate(d.position.X, d.position.Y)

	screen.DrawImage(d.sprite, op)
}

func (d *Droplet) Collider() Rect {
	bounds := d.sprite.Bounds()

	return NewRect(
		d.position.X,
		d.position.Y,
		float64(bounds.Dx()),
		float64(bounds.Dy()),
	)
}
