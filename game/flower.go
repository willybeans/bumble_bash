package game

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/willybeans/bumble_bash/assets"
)

const (
	rotationSpeedMin = -0.02
	rotationSpeedMax = 0.02
)

type Flower struct {
	position      Vector
	rotationSpeed float64
	sprite        *ebiten.Image
}

func NewFlower(baseVelocity float64) *Flower {
	pos := Vector{
		X: 0 + rand.Float64()*(screenWidth-50),
		Y: 0 + rand.Float64()*(screenHeight-50),
	}

	sprite := assets.FlowerSprite

	r := &Flower{
		position: pos,
		sprite:   sprite,
	}
	return r
}

func (f *Flower) Update() {
	// f.position.X += f.movement.X
	// f.position.Y += f.movement.Y
	// f.rotation += f.rotationSpeed
}

func (f *Flower) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(f.position.X, f.position.Y)

	screen.DrawImage(f.sprite, op)
}

func (f *Flower) Collider() Rect {
	bounds := f.sprite.Bounds()

	return NewRect(
		f.position.X,
		f.position.Y,
		float64(bounds.Dx()),
		float64(bounds.Dy()),
	)
}
