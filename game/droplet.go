package game

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/willybeans/bumble_bash/assets"
)

const (
	baseDropVelocity   = 5
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

	fmt.Println("target", target, shotCounter)

	// angle := rand.Float64() * 2 * math.Pi
	// radius := screenWidth / 2.0

	// fmt.Println("angle/radius", angle, radius)

	// pos := Vector{
	// 	X: target.X + math.Cos(angle)*radius,
	// 	Y: target.Y + math.Sin(angle)*radius,
	// }

	// fmt.Println("pos", pos)

	// velocity := baseDropVelocity + rand.Float64()*1.5
	velocity := float64(baseDropVelocity)

	fmt.Println("velocity", velocity)

	direction := Vector{
		X: target.X - spawnPos.X,
		Y: target.Y - spawnPos.Y,
	}

	// fmt.Println("direction", direction)

	normalizedDirection := direction.Normalize()

	movement := Vector{
		X: normalizedDirection.X * velocity,
		Y: normalizedDirection.Y * velocity,
	}

	fmt.Println("movement", movement)

	// sprite := assets.DropletSprite

	// d := &Droplet{
	// 	position:      spawnPos,
	// 	movement:      movement,
	// 	rotationSpeed: rotationSpeedMin + rand.Float64()*(rotationSpeedMax-rotationSpeedMin),
	// 	sprite:        sprite,
	// }

	// fmt.Println("droplet", d)

	// return d

	sprite := assets.DropletSprite

	// bounds := sprite.Bounds()
	// halfW := float64(bounds.Dx()) / 2
	// halfH := float64(bounds.Dy()) / 2

	// spawnPos.X -= halfW
	// spawnPos.Y -= halfH

	d := &Droplet{
		position: spawnPos,
		rotation: rotation,
		movement: movement,
		sprite:   sprite,
	}

	return d
}

func (d *Droplet) Update() {
	// 2 3 9 8 15 14
	// if d.rotation == 2 ||
	// 	d.rotation == 3 ||
	// 	d.rotation == 9 ||
	// 	d.rotation == 8 ||
	// 	d.rotation == 15 ||
	// 	d.rotation == 14 {
	// 	speed := dropSpeedPerSecond / float64(ebiten.TPS())
	// 	d.position.X += math.Cos(d.rotation) * speed
	// 	d.position.Y += math.Sin(d.rotation) * -speed
	// fmt.Println("test", d.rotation, d.position)

	//  3 {593.8003501586463 513.4143996238727}
	//  9 {597.4805877787148 500.7678040220514}
	// 15 {604.5478973999216 489.65323412600117}
	//  2 {620.5798142944668 477.5661200814684}
	//  8 {633.209998422265 473.82994849090886}
	// 14 {646.3810701830325 473.7716567342395}
	// }

	// d.position.X += d.rotation / 2
	// d.position.Y += d.rotation / 2
	d.position.X += d.movement.X
	d.position.Y += d.movement.Y
	// d.rotation += d.rotationSpeed
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
