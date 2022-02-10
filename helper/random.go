package helper

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync/atomic"
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

var num int64

func GenerateOrderNumber(t time.Time) string {
	s := t.Format(Continuity)
	m := t.UnixNano()/1e6 - t.UnixNano()/1e9*1e3
	ms := sup(m, 3)
	p := os.Getpid() % 1000
	ps := sup(int64(p), 3)
	i := atomic.AddInt64(&num, 1)
	r := i % 10000
	rs := sup(r, 4)
	n := fmt.Sprintf("%s%s%s%s", s, ms, ps, rs)
	return n
}

func sup(i int64, n int) string {
	m := fmt.Sprintf("%d", i)
	for len(m) < n {
		m = fmt.Sprintf("0%s", m)
	}
	return m
}
