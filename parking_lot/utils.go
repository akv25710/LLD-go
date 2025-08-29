package parking_lot

import (
	"math/rand"
	"time"
)

const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandomAlphaNumeric(n int) string {
	rand.Seed(time.Now().UnixNano())
	result := make([]byte, n)
	for i := 0; i < n; i++ {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}
