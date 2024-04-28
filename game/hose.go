package game

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/willybeans/bumble_bash/assets"
)

const (
	shootCooldown      = time.Millisecond * 500
	dropletSpawnOffset = 10.0
)

type Hose struct {
	game          *Game
	position      Vector
	sprite        *ebiten.Image
	rotation      float64
	hitCoords     Vector
	shootCooldown *Timer
}

func NewHose(game *Game) *Hose {
	sprite := assets.HoseSprite

	bounds := sprite.Bounds()
	// halfW := float64(bounds.Dx()) / 2
	// halfH := float64(bounds.Dy()) / 2

	// pos := Vector{
	// 	X: screenWidth/2 - halfW,
	// 	Y: screenHeight/2 - halfH,
	// }

	pos := Vector{
		X: screenWidth*.75 - float64(bounds.Dx())/2,
		Y: screenHeight*.75 - float64(bounds.Dy())/2,
	}

	return &Hose{
		game:          game,
		position:      pos,
		rotation:      0,
		sprite:        sprite,
		shootCooldown: NewTimer(shootCooldown)}
}

func (h *Hose) Update() {
	h.shootCooldown.Update()
	if h.shootCooldown.IsReady() {
		h.shootCooldown.Reset()

		spawnPos := Vector{
			X: h.position.X,
			Y: h.position.Y + 20,
		}

		droplet := NewDroplet(spawnPos, h.rotation)
		h.game.AddDroplet(droplet)
	}
}

func (h *Hose) Draw(screen *ebiten.Image) {
	var hoseSprite = h.sprite

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(h.position.X, h.position.Y)
	screen.DrawImage(hoseSprite, op)
}

// func (h *Hose) Collider() Rect {
// 	bounds := h.sprite.Bounds()
// 	return NewRect(
// 		h.position.X,
// 		h.position.Y,
// 		float64(bounds.Dx()),
// 		float64(bounds.Dy()),
// 	)
// }
