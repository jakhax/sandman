package utils;

import (
    "golang.org/x/sys/unix"
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

// Writable check if path is writable
func Writable(path string)(ok bool){
   if  err := unix.Access(path, unix.W_OK); err != nil{
       ok = false;
   }
    return;
}