package game

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/willybeans/bumble_bash/assets"
)

type Pollen struct {
	shouldFall      bool
	position        Vector
	pollenPlacement Vector
	rotation        float64
	movement        Vector
	rotationSpeed   float64
	sprite          *ebiten.Image
}

func NewPollen(baseVelocity float64, playerSprite *ebiten.Image) *Pollen {
	var (
		innerW = float64(playerSprite.Bounds().Dx()) * float64(0.75)
		innerH = float64(playerSprite.Bounds().Dy()) * float64(0.75)
		height = float64(playerSprite.Bounds().Dy())
		width  = float64(playerSprite.Bounds().Dx())
	)

	innerBox := Vector{
		X: (width - innerW) / 2,
		Y: (height - innerH) / 2,
	}

	position := Vector{
		X: randFloatRange(innerBox.X, innerW),
		Y: randFloatRange(innerBox.Y, innerH),
	}

	velocity := baseVelocity + 5

	// direction := Vector{
	// 	X: target.X - pos.X,
	// 	Y: target.Y - pos.Y,
	// }
	direction := Vector{
		X: 0,
		Y: 0,
	}
	normalizedDirection := direction.Normalize()

	movement := Vector{
		X: normalizedDirection.X * velocity,
		Y: normalizedDirection.Y * velocity,
	}

	sprite := assets.PollenSprites[rand.Intn(len(assets.PollenSprites))]

	r := &Pollen{
		shouldFall:      false,
		position:        position,
		pollenPlacement: innerBox,
		movement:        movement,
		sprite:          sprite,
	}
	return r
}

func (p *Pollen) Update() {
	if p.shouldFall {
		p.position.X += p.movement.X
		p.position.Y += p.movement.Y
		p.rotation += p.rotationSpeed
	}

}

func (p *Pollen) Draw(screen *ebiten.Image, pos Vector) {

	if p.shouldFall {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Rotate(p.rotation)
		op.GeoM.Translate(pos.X, pos.Y)
		screen.DrawImage(p.sprite, op)
	}
}

func (p *Pollen) Collider() Rect {
	bounds := p.sprite.Bounds()

	return NewRect(
		p.position.X,
		p.position.Y,
		float64(bounds.Dx()),
		float64(bounds.Dy()),
	)
}
