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
