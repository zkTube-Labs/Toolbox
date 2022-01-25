package helper

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

var defaultLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// RandomString returns a random string with a fixed length
func RandomString(n int, allowedChars ...[]rune) string {
	rand.Seed(time.Now().UnixNano() + int64(rand.Intn(100)))
	var letters []rune
	if len(allowedChars) == 0 {
		letters = defaultLetters
	} else {
		letters = allowedChars[0]
	}

	b := make([]rune, n)

	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

func RandomNumber(n int) string {
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	str := "9"
	for i := 0; i < n-1; i++ {
		str += "9"
	}
	i, _ := strconv.Atoi(str)
	randomNumber := r.Intn(i)
	return fmt.Sprintf("%0"+strconv.Itoa(n)+"d", randomNumber)
}
