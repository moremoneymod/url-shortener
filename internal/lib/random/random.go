package random

import (
	"math/rand"
	"time"
)

func NewRandomString(stringLength int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghijklmnopqrstuvwxyz" + "0123456789" + "!@#$%^&*()_+=")

	b := make([]rune, stringLength)
	for i := range b {
		b[i] = chars[rnd.Intn(len(chars))]
	}
	return string(b)
}
