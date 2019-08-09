package compute

import (
	"math/rand"
	"time"
)

func testRandomString() string {
	chars := "1234567890abcdefghijklmnopqrstuvwxyz"

	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 10)
	for i := range b {
		b[i] = chars[rand.Int63()%int64(len(chars))]
	}

	return string(b)
}
