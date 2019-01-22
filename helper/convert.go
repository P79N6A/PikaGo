package helper

import "strconv"

func S2I(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func I2S(i int) string {
	return strconv.Itoa(i)
}
