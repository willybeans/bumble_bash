package game

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/willybeans/bumble_bash/assets"
)

const (
	apiarySpawnTime = time.Millisecond * 4000
)

type Apiary struct {
	game             *Game
	didSpawn         bool
	position         Vector
	sprite           *ebiten.Image
	apiarySpawnTimer *Timer
}

func NewApiary(game *Game) *Apiary {
	sprite := assets.ApiarySprite

	bounds := sprite.Bounds()

	pos := Vector{
		X: screenWidth - float64(bounds.Dx()*4),
		Y: 0 + float64(bounds.Dy()),
	}

	return &Apiary{
		game:             game,
		didSpawn:         false,
		position:         pos,
		sprite:           sprite,
		apiarySpawnTimer: NewTimer(apiarySpawnTime)}
}

func (a *Apiary) Update() {
	if a.didSpawn == false {
		a.apiarySpawnTimer.Update()
		if a.apiarySpawnTimer.IsReady() {
			a.didSpawn = true
			a.apiarySpawnTimer.Reset()
		}
	}
}

func (a *Apiary) Draw(screen *ebiten.Image) {
	var apiarySprite = a.sprite

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(4, 4)
	op.GeoM.Translate(a.position.X, a.position.Y)
	screen.DrawImage(apiarySprite, op)
}

func (a *Apiary) Collider() Rect {
	bounds := a.sprite.Bounds()
	return NewRect(
		a.position.X,
		a.position.Y,
		float64(bounds.Dx()*4),
		float64(bounds.Dy()*4),
	)
}
