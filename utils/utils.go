package utils

import "strings"

const PaddingLeft = "left"
const PaddingRight = "right"

func Padding(input string, length int, char string, direction string) string {
	if len(input) >= length {
		return input
	}

	if direction == "left" {
		return strings.Repeat(char, length-len(input)) + input
	} else if direction == "right" {
		return input + strings.Repeat(char, length-len(input))
	} else {
		return input
	}
}

func Modulo(dividend string, divisor int) int {
	var remainder int

	for i := 0; i < len(dividend); i++ {
		remainder = (remainder*10 + int(dividend[i]-'0')) % divisor
	}

	return remainder
}
