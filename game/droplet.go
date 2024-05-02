package game

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/willybeans/bumble_bash/assets"
)

const (
	baseDropVelocity   = 3
	dropSpeedPerSecond = 350.0
)

type Droplet struct {
	position      Vector
	rotation      float64
	movement      Vector
	rotationSpeed float64
	sprite        *ebiten.Image
}

func NewDroplet(spawnPos Vector, rotation float64, shotCounter int) *Droplet {
	var target Vector
	if shotCounter < 5 {
		target = Vector{
			X: rotation * 100,
			Y: screenHeight - float64(shotCounter)*100,
		}
	} else if shotCounter == 5 {
		target = Vector{
			X: rotation * 100,
			Y: rotation * 100,
		}
	} else {
		target = Vector{
			X: 0 + float64(shotCounter)*50,
			Y: rotation * 100,
		}
	}

	velocity := float64(baseDropVelocity)

	direction := Vector{
		X: target.X - spawnPos.X,
		Y: target.Y - spawnPos.Y,
	}

	normalizedDirection := direction.Normalize()

	movement := Vector{
		X: normalizedDirection.X * velocity,
		Y: normalizedDirection.Y * velocity,
	}

	sprite := assets.DropletSprite

	d := &Droplet{
		position: spawnPos,
		rotation: rotation,
		movement: movement,
		sprite:   sprite,
	}

	return d
}

func (d *Droplet) Update() {
	d.position.X += d.movement.X
	d.position.Y += d.movement.Y
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
