package utils

import (
	"unsafe"
	"math/rand"
	"time"
)

const (
	source        = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789+_"
	sourceLength  = len(source)
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)
var src = rand.NewSource(time.Now().UnixNano())


// Generates a random salt of n characters.
func Salt(n int) string {
	b := make([]byte, n)
	// This loop goes like this:
	// While there are still letters to fill, get a random number of
	//   63 bits, and consume it in chunks of 6 bits to take a letter
	//   from the 64-size index source (64 slots = 6 bits).
	//   If there are not enough bits to consume, get a new random
	//   number and continue until all the slots are filled.
	//
	// Inspired on: https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
	//   with slight changes.
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain < letterIdxBits {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < sourceLength {
			b[i] = source[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}
