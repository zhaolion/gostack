package stringutil

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// UniqueID returns randomly generated string with prefix
func UniqueID(prefix string, l uint) string {
	return fmt.Sprintf("%s%s", prefix, RandomString(l))
}

func UniqueIDBaseTime(prefix string) string {
	timeStr := strconv.FormatUint(uint64(time.Now().UnixNano()), 36)
	return fmt.Sprintf("%s%s%s", prefix, timeStr, RandomString(2))
}

// RandomString returns randomly generated string
func RandomString(l uint) string {
	s := make([]byte, l)
	for i := 0; i < int(l); i++ {
		s[i] = chars[rand.Intn(len(chars))]
	}
	return string(s)
}
