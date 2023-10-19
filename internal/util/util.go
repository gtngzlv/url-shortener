package util

import (
	"crypto/rand"
	"math/big"
)

const codeAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const defaultLength = 10

var codeAlphabetLen = big.NewInt(int64(len(codeAlphabet)))

// RandStringRunes return random strings for short url generation
func RandStringRunes() string {
	b := make([]byte, defaultLength)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	for i, v := range b {
		b[i] = codeAlphabet[int(v)%len(codeAlphabet)]
	}
	return string(b)
}
