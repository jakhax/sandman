package utils;

import (
	"time"
	"math/rand"
)


// RandString generates pseudo random string of length n
func RandString(n int) string {
	rand.Seed(time.Now().UnixNano())
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ";
    b := make([]byte, n);
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))];
    }
    return string(b);
}