package yacli

import (
	"crypto/rand"
	"fmt"
)

func uniqId() string {
	n := 5
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%X", b)
}
