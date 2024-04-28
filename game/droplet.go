package game

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/willybeans/bumble_bash/assets"
)

const (
	baseDropVelocity = 5
)

type Droplet struct {
	position      Vector
	rotation      float64
	movement      Vector
	rotationSpeed float64
	sprite        *ebiten.Image
}

func NewDroplet(spawnPos Vector, rot float64) *Droplet {
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

	velocity := baseDropVelocity + rand.Float64()*1.5

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
		position:      spawnPos,
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
	op := &ebiten.DrawImageOptions{}
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
