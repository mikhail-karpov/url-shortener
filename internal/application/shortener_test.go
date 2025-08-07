package application

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShorten(t *testing.T) {

	t.Run("shortens to alphanumeric string", func(t *testing.T) {

		type testCase struct {
			id       uint32
			expected string
		}

		testCases := []testCase{
			{0, "0"},
			{1024, "GW"},
			{10241024, "gy9o"},
		}

		for _, c := range testCases {
			actual := shorten(c.id)
			assert.Equal(t, c.expected, actual)
			assert.LessOrEqual(t, len(actual), 7)
		}
	})

	t.Run("is idempotent", func(t *testing.T) {

		for i := 0; i < 1000; i++ {
			actual := shorten(10241024)
			assert.Equal(t, "gy9o", actual)
		}
	})
}
