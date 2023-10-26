package util_test

import (
	"fmt"

	"github.com/gtngzlv/url-shortener/internal/util"
)

func ExampleRandStringRunes() {
	shortID := util.RandStringRunes()
	shortID = "some id like kjfhsl"
	fmt.Println(shortID)
	// Output:
	// some id like kjfhsl
}
