package slugs

import (
	"math/rand"
	"strings"
	"time"
)

var (
	chars    []string
	numChars int
)

func init() {
	chars = []string{"A", "a", "B", "b", "C", "c", "D", "d", "E", "e", "F", "f", "G", "g", "H", "h", "I", "i", "J", "j", "K", "k", "L", "l", "M", "m", "N", "n", "O", "o", "P", "p", "Q", "q", "R", "r", "S", "s", "T", "t", "U", "u", "V", "v", "W", "w", "X", "x", "Y", "y", "Z", "z", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	numChars = len(chars)
	rand.Seed(time.Now().UnixNano())
}

// Random returns a random slug of length n
func Random(n int) string {
	var b strings.Builder
	for i := 0; i < n; i++ {
		b.WriteString(chars[rand.Intn(numChars)])
	}
	return b.String()
}
