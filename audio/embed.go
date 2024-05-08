package audio

import (
	_ "embed"
)

var (
	//go:embed ouch0.wav
	Ouch0_wav []byte

	//go:embed ouch1.wav
	Ouch1_wav []byte

	//go:embed ouch2.wav
	Ouch2_wav []byte

	//go:embed ouch2.wav
	Ouch3_wav []byte

	//go:embed stayup.wav
	Stayup_wav []byte

	// //go:embed ragtime.ogg
	// Ragtime_ogg []byte
)

var Ouchies = [][]byte{Ouch0_wav, Ouch1_wav, Ouch2_wav, Ouch3_wav}
