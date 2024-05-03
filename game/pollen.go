package game

import (
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/willybeans/bumble_bash/assets"
)

const (
	fallCooldown       = time.Millisecond * 1000
	basePollenVelocity = 2.0
)

type Pollen struct {
	shouldFall      bool
	catchable       bool
	position        Vector
	pollenPlacement Vector
	rotation        float64
	movement        Vector
	rotationSpeed   float64
	fallCooldown    *Timer
	sprite          *ebiten.Image
}

func NewPollen(baseVelocity float64, playerSprite *ebiten.Image, playerPosition Vector, shouldFall bool) *Pollen {
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

	velocity := basePollenVelocity

	direction := Vector{
		X: 0,
		Y: 0,
	}
	if shouldFall {
		direction = Vector{
			X: playerPosition.X,
			Y: screenHeight,
		}
	}
	normalizedDirection := direction.Normalize()

	movement := Vector{
		X: normalizedDirection.X * velocity,
		Y: normalizedDirection.Y * velocity,
	}

	sprite := assets.PollenSprites[rand.Intn(len(assets.PollenSprites))]

	if shouldFall {
		fallPosition := Vector{
			X: playerPosition.X + float64(playerSprite.Bounds().Dx()/2),
			Y: playerPosition.Y + float64(playerSprite.Bounds().Dy()/2),
		}
		position = fallPosition
	}

	r := &Pollen{
		shouldFall:      shouldFall,
		catchable:       false,
		position:        position,
		pollenPlacement: innerBox,
		movement:        movement,
		sprite:          sprite,
		fallCooldown:    NewTimer(fallCooldown),
	}
	return r
}

func (p *Pollen) Update() {
	if p.shouldFall {
		p.position.X += p.movement.X
		p.position.Y += p.movement.Y

		p.fallCooldown.Update()
		if p.fallCooldown.IsReady() {
			if p.catchable == false {
				p.catchable = true
			}
			p.fallCooldown.Reset()
		}
	}
}

func (p *Pollen) Draw(screen *ebiten.Image) {
	if p.shouldFall {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Rotate(p.rotation)
		op.GeoM.Translate(p.position.X, p.position.Y)
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
