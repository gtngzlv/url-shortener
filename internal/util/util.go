package util

import (
	"crypto/rand"
	"log"
	"math/big"
)

var defaultLength = 8

func RandStringRunes() string {
	var codeAlphabet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := ""
	for i := 0; i < defaultLength; i++ {
		b += string(codeAlphabet[cryptoRandSecure(int64(len(codeAlphabet)))])
	}
	return b
}

func cryptoRandSecure(max int64) int64 {
	nBig, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		log.Println(err)
	}
	return nBig.Int64()
}
