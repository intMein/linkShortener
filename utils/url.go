package utils

import (
	"bytes"
	"math/rand"
)

var urlAlphabet = []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}

func GenURL(n int) string {
	buf := bytes.NewBufferString("")

	for i := 0; i < n; i++ {
		buf.WriteByte(urlAlphabet[rand.Intn(len(urlAlphabet))])
	}

	return buf.String()
}
