package number

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUtils(t *testing.T) {
	t.Run("FormatNumber", func(t *testing.T) {
		t.Run("should format number correctly", func(t *testing.T) {
			result := FormatNumber("+1 555-444-2222")

			assert.Equal(t, result, "15554442222", "they should be equal")
		})

		t.Run("should format number correctly", func(t *testing.T) {
			result := FormatNumber("+1 (315) 284-1580")

			assert.Equal(t, result, "13152841580", "they should be equal")
		})
	})

	t.Run("IsValid", func(t *testing.T) {
		t.Run("should validate phone number", func(t *testing.T) {
			result := IsValid("+1 315-284-1580")

			assert.Equal(t, result, true, "they should be equal")
		})

		t.Run("should validate phone number", func(t *testing.T) {
			result := IsValid("P+1 315-284-1580A")

			assert.Equal(t, result, false, "they should be equal")
		})
	})
}
