package math64

import "math/rand"

func RandInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}
