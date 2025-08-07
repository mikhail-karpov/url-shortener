package application

import (
	"strings"
)

func shorten(id uint32) string {
	var indexes []uint32
	if id == 0 {
		indexes = []uint32{0}
	} else {
		indexes = randomIndexes(id)
	}

	builder := strings.Builder{}
	for _, i := range indexes {
		builder.WriteString(string(alphabet[i]))
	}
	return builder.String()
}

func randomIndexes(id uint32) []uint32 {
	num := id
	digits := []uint32{}

	for num > 0 {
		digits = append(digits, num%alphabetLen)
		num /= alphabetLen
	}

	// reverse digits
	for i, j := 0, len(digits)-1; i < j; i, j = i+1, j-1 {
		digits[i], digits[j] = digits[j], digits[i]
	}
	return digits
}
