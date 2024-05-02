package game

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/willybeans/bumble_bash/assets"
)

const (
	shootCooldown      = time.Millisecond * 2000
	dropletSpawnOffset = 10.0
)

type Hose struct {
	id            int
	game          *Game
	position      Vector
	sprite        *ebiten.Image
	rotation      float64
	hitCoords     Vector
	shootCooldown *Timer
	shotCounter   int
}

func NewHose(game *Game, hoseCount int) *Hose {
	sprite := assets.HoseSprite

	bounds := sprite.Bounds()

	var pos Vector
	if hoseCount%2 == 0 {
		pos = Vector{
			X: 0 + float64(bounds.Dx()/10),
			Y: screenHeight - float64(bounds.Dy()),
		}
	} else {
		pos = Vector{
			X: screenWidth - float64(bounds.Dx()),
			Y: screenHeight - float64(bounds.Dy()),
		}
	}

	// pos := Vector{
	// 	X: screenWidth - float64(bounds.Dx()),
	// 	Y: screenHeight - float64(bounds.Dy()),
	// }

	return &Hose{
		id:            hoseCount,
		game:          game,
		position:      pos,
		rotation:      0,
		sprite:        sprite,
		shotCounter:   1,
		shootCooldown: NewTimer(shootCooldown)}
}

func (h *Hose) Update() {

	h.shootCooldown.Update()
	if h.shootCooldown.IsReady() {
		h.shotCounter++
		h.shootCooldown.Reset()

		spawnPos := Vector{
			X: h.position.X,
			Y: h.position.Y + float64(h.sprite.Bounds().Dy()/2),
		}

		for i := 1; i < 6; i++ {
			h.rotation = float64(i)
			// fmt.Println("counter", h.shotCounter)
			// fmt.Println(float64(h.shotCounter % 5))
			droplet := NewDroplet(spawnPos, h.rotation, h.shotCounter)
			h.game.AddDroplet(droplet)
		}
		if h.shotCounter >= 10 {
			h.shotCounter = 1
		}

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
