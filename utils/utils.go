package utils

import "math/rand"

func RandString(n int) string {
	lenght := rand.Int()%n + 1

	str := make([]byte, lenght)
	for i := 0; i < lenght; i++ {
		str[i] = byte(rand.Int()%26) + 'a'
	}
	return string(str)
}
