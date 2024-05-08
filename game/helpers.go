package game

import (
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

func randFloatRange(min, max float64) float64 {
	return rand.Float64() * (max - min)
}

func randIntRange(min, max int) int {
	return rand.Intn(max-min) + min
}

type AnyType interface {
	Flower | Droplet | Hose | Pollen
}

func RemoveIndex[T AnyType](slc []*T, index int, remove *T) []*T {
	for idx, v := range slc {
		if v == remove {
			return append(slc[0:idx], slc[idx+1:]...)
		}
	}
	return slc
}

func playSound(sound *audio.Player) {
	if sound.IsPlaying() {
		return
	} else {
		sound.Rewind()
		sound.Play()
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
