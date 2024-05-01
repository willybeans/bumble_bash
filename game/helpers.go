package game

import (
	"math/rand"
)

func randFloatRange(min, max float64) float64 {
	return rand.Float64() * (max - min)
}

func randIntRange(min, max int) int {
	return rand.Intn(max-min) + min
}

type AnyType interface{ Flower | Droplet | Hose }

func RemoveIndex[T AnyType](slc []*T, index int, remove *T) []*T {
	for idx, v := range slc {
		if v == remove {
			return append(slc[0:idx], slc[idx+1:]...)
		}
	}
	return slc
}
