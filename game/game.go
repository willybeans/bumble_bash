package game

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	raudio "github.com/willybeans/bumble_bash/audio"
)

type Mode int

const (
	ModeTitle Mode = iota
	ModeGame
	ModeGameOver

	sampleRate = 48000

	screenWidth   = 800
	screenHeight  = 600
	tileSize      = 32
	titleFontSize = fontSize * 1.5
	fontSize      = 24
	smallFontSize = fontSize / 2

	flowerSpawnTime  = 1 * time.Second
	hoseSpawnTime    = 12 * time.Second
	dropletSpawnTime = 2 * time.Second

	baseFlowerVelocity  = 0.25
	flowerSpeedUpAmount = 0.1
	flowerSpeedUpTime   = 5 * time.Second

	introLengthInSecond = 76
	loopLengthInSecond  = 4
	bytesPerSample      = 4 // 2 channels * 2 bytes (16 bit)
)

var (
	isHit            = false
	arcadeFaceSource *text.GoTextFaceSource
)

type Game struct {
	mode     Mode
	player   *Player
	apiary   *Apiary
	flowers  []*Flower
	hoses    []*Hose
	droplets []*Droplet
	pollens  []*Pollen

	flowerSpawnTimer  *Timer
	hoseSpawnTimer    *Timer
	dropletSpawnTimer *Timer

	score int

	baseVelocity  float64
	velocityTimer *Timer

	audioContext *audio.Context
	audioPlayer  *audio.Player
	hitPlayer    []*audio.Player
	hitPollen    *audio.Player
	hitApiary    *audio.Player
	hitFlower    *audio.Player
	songStarted  bool
}

func NewGame() *Game {
	g := &Game{
		dropletSpawnTimer: NewTimer(dropletSpawnTime),
		flowerSpawnTimer:  NewTimer(flowerSpawnTime),
		hoseSpawnTimer:    NewTimer(hoseSpawnTime),
		baseVelocity:      baseFlowerVelocity,
		velocityTimer:     NewTimer(flowerSpeedUpTime),
	}

	g.player = NewPlayer(g)
	g.apiary = NewApiary(g)

	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.PressStart2P_ttf))
	if err != nil {
		log.Fatal(err)
	}
	arcadeFaceSource = s

	if g.audioContext == nil {
		g.audioContext = audio.NewContext(sampleRate)
	}

	for i := 0; i < len(raudio.Ouchies); i++ {
		hitD, err := wav.DecodeWithoutResampling(bytes.NewReader(raudio.Ouchies[i]))
		check(err)
		hitPlayer, err := g.audioContext.NewPlayer(hitD)
		hitPlayer.SetVolume(0.5)
		g.hitPlayer = append(g.hitPlayer, hitPlayer)
		check(err)
	}

	wowS, err := wav.DecodeWithoutResampling(bytes.NewReader(raudio.Wow_wav))
	check(err)
	g.hitApiary, err = g.audioContext.NewPlayer(wowS)
	g.hitApiary.SetVolume(1)
	check(err)

	NomS, err := wav.DecodeWithoutResampling(bytes.NewReader(raudio.Nom_wav))
	check(err)
	g.hitFlower, err = g.audioContext.NewPlayer(NomS)
	g.hitFlower.SetVolume(0.1)
	check(err)

	GrabS, err := wav.DecodeWithoutResampling(bytes.NewReader(raudio.Grab_wav))
	check(err)
	g.hitPollen, err = g.audioContext.NewPlayer(GrabS)
	check(err)

	return g
}
func (g *Game) isKeyJustPressed() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		return true
	}
	return false
}

func (g *Game) Update() error {
	switch g.mode {
	case ModeTitle:
		if g.isKeyJustPressed() {
			g.mode = ModeGame
		}
	case ModeGame:
		// g.x16 += 32
		// g.cameraX += 2
		// if g.isKeyJustPressed() {
		// 	g.vy16 = -96
		// 	if err := g.jumpPlayer.Rewind(); err != nil {
		// 		return err
		// 	}
		// 	g.jumpPlayer.Play()
		// }
		// g.y16 += g.vy16

		// // Gravity
		// g.vy16 += 4
		// if g.vy16 > 96 {
		// 	g.vy16 = 96
		// }

		// if g.hit() {
		// 	if err := g.hitPlayer.Rewind(); err != nil {
		// 		return err
		// 	}
		// 	g.hitPlayer.Play()
		// 	g.mode = ModeGameOver
		// 	g.gameoverCount = 30
		// }
	case ModeGameOver:
		// if g.gameoverCount > 0 {
		// 	g.gameoverCount--
		// }
		// if g.gameoverCount == 0 && g.isKeyJustPressed() {
		// 	g.init()
		// 	g.mode = ModeTitle
		// }
	}
	g.velocityTimer.Update()
	if g.velocityTimer.IsReady() {
		g.velocityTimer.Reset()
		g.baseVelocity += flowerSpeedUpAmount
	}

	g.player.Update()
	if g.apiary.didSpawn { // move the timer to here for the spawning
	}
	g.apiary.Update()

	if g.apiary.Collider().Intersects(g.player.Collider()) {
		if len(g.player.Pollens) > 0 {
			score := len(g.player.Pollens)
			g.player.Pollens = nil
			g.score += score
			playSound(g.hitApiary)
		}

	}

	g.flowerSpawnTimer.Update()
	if g.flowerSpawnTimer.IsReady() {
		g.flowerSpawnTimer.Reset()

		f := NewFlower(g.baseVelocity)
		g.flowers = append(g.flowers, f)
	}

	g.hoseSpawnTimer.Update()
	if g.hoseSpawnTimer.IsReady() {
		g.hoseSpawnTimer.Reset()

		if len(g.hoses) < 2 {
			h := NewHose(g, len(g.hoses))
			g.hoses = append(g.hoses, h)
		}
	}

	for _, f := range g.flowers {
		f.Update()
	}

	for _, d := range g.droplets {
		d.Update()
	}

	for _, p := range g.pollens {
		p.Update()
	}

	for _, h := range g.hoses {
		h.Update()
	}

	for i, f := range g.flowers {
		isPlayerCollision := f.Collider().Intersects(g.player.Collider())

		if f.position.X < 0 ||
			f.position.Y < 0 ||
			f.position.X > screenWidth ||
			f.position.Y > screenHeight ||
			isPlayerCollision {

			g.flowers = RemoveIndex(g.flowers, i, f)

			if isPlayerCollision {
				p := NewPollen(g.baseVelocity, g.player.sprite, g.player.position, false)
				g.player.Pollens = append(g.player.Pollens, p)
				playSound(g.hitFlower)

			}
		}

	}

	for i, d := range g.droplets {
		isPlayerCollision := d.Collider().Intersects(g.player.Collider())
		if d.position.X < 0 ||
			d.position.Y < 0 ||
			d.position.X > screenWidth ||
			d.position.Y > screenHeight ||
			isPlayerCollision {

			g.droplets = RemoveIndex(g.droplets, i, d)

			if isPlayerCollision {
				g.player.isHit = true
				g.player.hitCoords = Vector{
					X: d.Collider().X,
					Y: d.Collider().Y,
				}

				randomOuch := randIntRange(0, len(raudio.Ouchies))
				playSound(g.hitPlayer[randomOuch])

				if len(g.player.Pollens) > 0 {
					g.player.Pollens = g.player.Pollens[:len(g.player.Pollens)-1]
					pollen := NewPollen(g.baseVelocity, g.player.sprite, g.player.position, true)
					g.pollens = append(g.pollens, pollen)
				}
			}

		}
	}

	for i, p := range g.player.Pollens {
		if p.position.X < 0 ||
			p.position.Y < 0 ||
			p.position.X > screenWidth ||
			p.position.Y > screenHeight {

			g.player.Pollens = RemoveIndex(g.player.Pollens, i, p)

		}
	}

	for i, p := range g.pollens {
		isPlayerCollision := p.Collider().Intersects(g.player.Collider())

		if p.position.X < 0 ||
			p.position.Y < 0 ||
			p.position.X > screenWidth ||
			p.position.Y > screenHeight ||
			(isPlayerCollision && p.catchable) {

			if isPlayerCollision && p.catchable {
				pollen := NewPollen(g.baseVelocity, g.player.sprite, g.player.position, false)
				g.player.Pollens = append(g.player.Pollens, pollen)
				playSound(g.hitPollen)
			}
			g.pollens = RemoveIndex(g.pollens, i, p)

		}
	}

	if g.songStarted == false {
		// g.audioContext = audio.NewContext(sampleRate)
		// Decode an Ogg file.
		// oggS is a decoded io.ReadCloser and io.Seeker.
		oggS, err := wav.DecodeWithoutResampling(bytes.NewReader(raudio.Stayup_wav))
		if err != nil {
			return err
		}

		// Create an infinite loop stream from the decoded bytes.
		// s is still an io.ReadCloser and io.Seeker.
		s := audio.NewInfiniteLoop(oggS, introLengthInSecond*bytesPerSample*sampleRate)

		g.audioPlayer, err = g.audioContext.NewPlayer(s)
		if err != nil {
			return err
		}

		// Play the infinite-length stream. This never ends.
		g.audioPlayer.SetVolume(0.1)
		g.audioPlayer.Play()
		g.songStarted = true
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x80, 0xa0, 0xc0, 0xff})

	var titleTexts string
	var texts string
	switch g.mode {
	case ModeTitle:
		titleTexts = "LET'S GET\n READY TO BUMBLE"
		texts = "\n\n\n\n\n\nPRESS SPACE KEY\n"
	case ModeGameOver:
		texts = "\nGAME OVER!"
	}

	op := &text.DrawOptions{}
	op.GeoM.Translate(screenWidth/2, 3*titleFontSize)
	op.ColorScale.ScaleWithColor(color.White)
	op.LineSpacing = titleFontSize
	op.PrimaryAlign = text.AlignCenter
	text.Draw(screen, titleTexts, &text.GoTextFace{
		Source: arcadeFaceSource,
		Size:   titleFontSize,
	}, op)

	op = &text.DrawOptions{}
	op.GeoM.Translate(screenWidth/2, 3*titleFontSize)
	op.ColorScale.ScaleWithColor(color.White)
	op.LineSpacing = fontSize
	op.PrimaryAlign = text.AlignCenter
	text.Draw(screen, texts, &text.GoTextFace{
		Source: arcadeFaceSource,
		Size:   fontSize,
	}, op)

	if g.mode == ModeTitle {
		const msg = "Let's Get Ready To Bumble By Will Wedmedyk is\nlicenced under CC BY 3.0."

		op := &text.DrawOptions{}
		op.GeoM.Translate(screenWidth/2, screenHeight-smallFontSize/2)
		op.ColorScale.ScaleWithColor(color.White)
		op.LineSpacing = smallFontSize
		op.PrimaryAlign = text.AlignCenter
		op.SecondaryAlign = text.AlignEnd
		text.Draw(screen, msg, &text.GoTextFace{
			Source: arcadeFaceSource,
			Size:   smallFontSize,
		}, op)
	}

	op = &text.DrawOptions{}
	op.GeoM.Translate(screenWidth/2, 0)
	op.ColorScale.ScaleWithColor(color.White)
	op.LineSpacing = fontSize
	op.PrimaryAlign = text.AlignCenter
	text.Draw(screen, fmt.Sprintf("%04d", g.score), &text.GoTextFace{
		Source: arcadeFaceSource,
		Size:   fontSize,
	}, op)

	// acceleration linear drag -- physics -- motion of linear drag
	//

	if g.mode != ModeTitle {
		g.player.Draw(screen)

		g.apiary.Draw(screen)

		for _, f := range g.flowers {
			f.Draw(screen)
		}
		for _, h := range g.hoses {
			h.Draw(screen)
		}
		for _, d := range g.droplets {
			d.Draw(screen)
		}
		for _, p := range g.pollens {
			p.Draw(screen)
		}
		for _, p := range g.player.Pollens {
			p.Draw(screen)
		}

		// for _, f := range g.flowers {
		// 	vector.StrokeRect(
		// 		screen,
		// 		float32(f.position.X),
		// 		float32(f.position.Y),
		// 		float32(f.sprite.Bounds().Dx()),
		// 		float32(f.sprite.Bounds().Dy()),
		// 		1.0,
		// 		color.White,
		// 		false,
		// 	)
		// }

		// for _, d := range g.droplets {
		// 	vector.StrokeRect(
		// 		screen,
		// 		float32(d.position.X),
		// 		float32(d.position.Y),
		// 		float32(d.sprite.Bounds().Dx()),
		// 		float32(d.sprite.Bounds().Dy()),
		// 		1.0,
		// 		color.White,
		// 		false,
		// 	)
		// }

		// vector.StrokeRect(
		// 	screen,
		// 	float32(g.player.position.X),
		// 	float32(g.player.position.Y),
		// 	float32(g.player.sprite.Bounds().Dx()),
		// 	float32(g.player.sprite.Bounds().Dy()),
		// 	1.0,
		// 	color.White,
		// 	false,
		// )

		// vector.StrokeRect(
		// 	screen,
		// 	float32(g.apiary.position.X),
		// 	float32(g.apiary.position.Y),
		// 	float32(g.apiary.sprite.Bounds().Dx()*4),
		// 	float32(g.apiary.sprite.Bounds().Dy()*4),
		// 	1.0,
		// 	color.White,
		// 	false,
		// )

		// for _, h := range g.hoses {
		// 	vector.StrokeRect(
		// 		screen,
		// 		float32(h.position.X),
		// 		float32(h.position.Y),
		// 		float32(h.sprite.Bounds().Dx()),
		// 		float32(h.sprite.Bounds().Dy()),
		// 		1.0,
		// 		color.White,
		// 		false,
		// 	)
		// }

	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) AddDroplet(d *Droplet) {
	g.droplets = append(g.droplets, d)
}
