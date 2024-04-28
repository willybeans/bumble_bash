package game

import (
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/willybeans/bumble_bash/assets"
)

const (
	hitCooldown       = time.Millisecond * 500
	rotationPerSecond = math.Pi

	bulletSpawnOffset = 50.0
)

type Player struct {
	game     *Game
	position Vector
	// rotation float64
	oldPollenCount int
	Pollens        []*Pollen
	sprite         *ebiten.Image
	hitSprite      *ebiten.Image
	hitCooldown    *Timer
	isHit          bool
	hitCoords      Vector
}

func NewPlayer(game *Game) *Player {
	sprite := assets.PlayerSprite
	hitSprite := assets.PlayerHitSprite

	bounds := sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	pos := Vector{
		X: screenWidth/2 - halfW,
		Y: screenHeight/2 - halfH,
	}

	return &Player{
		game:           game,
		position:       pos,
		oldPollenCount: 0,
		sprite:         sprite,
		hitSprite:      hitSprite,
		hitCooldown:    NewTimer(hitCooldown),
	}
}

func (p *Player) Update() {
	playerSpeed := 6.0
	var delta Vector

	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		delta.Y = playerSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		delta.Y = -playerSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		delta.X = -playerSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		delta.X = playerSpeed
	}

	// Check for diagonal movement
	if delta.X != 0 && delta.Y != 0 {
		factor := playerSpeed / math.Sqrt(delta.X*delta.X+delta.Y*delta.Y)
		delta.X *= factor
		delta.Y *= factor
	}

	p.position.Y += delta.Y
	p.position.X += delta.X

	if p.isHit {
		p.hitCooldown.Update()
		if p.hitCooldown.IsReady() {
			p.hitCooldown.Reset()
			p.isHit = false
		}
	}

}

func (p *Player) Draw(screen *ebiten.Image) {
	var playerSprite = p.sprite
	if p.isHit {
		playerSprite = p.hitSprite
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.position.X, p.position.Y)
	screen.DrawImage(playerSprite, op)

	for _, pollen := range p.Pollens {
		opc := &ebiten.DrawImageOptions{}
		opc.GeoM.Translate(p.position.X, p.position.Y)
		opc.GeoM.Translate(pollen.pollenPlacement.X, pollen.pollenPlacement.Y)
		opc.GeoM.Translate(pollen.position.X, pollen.position.Y)
		screen.DrawImage(pollen.sprite, opc)
	}

}

func (p *Player) Collider() Rect {
	bounds := p.sprite.Bounds()
	return NewRect(
		p.position.X,
		p.position.Y,
		float64(bounds.Dx()),
		float64(bounds.Dy()),
	)
}
