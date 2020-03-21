package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtils(t *testing.T) {
	assert := assert.New(t)

	t.Run("FormatNumber", func(t *testing.T) {
		t.Run("should format number correctly", func(t *testing.T) {
			result := FormatNumber("+1 555-444-2222")

			assert.Equal(result, "15554442222", "they should be equal")
		})

		t.Run("should format number correctly", func(t *testing.T) {
			result := FormatNumber("+1 (315) 284-1580")

			assert.Equal(result, "13152841580", "they should be equal")
		})
	})

	t.Run("ParseCountryCode", func(t *testing.T) {
		t.Run("should parse country code correctly", func(t *testing.T) {
			result := ParseCountryCode("+33 679368229")

			assert.Equal(result, "FR", "they should be equal")
		})

		t.Run("should parse country code correctly", func(t *testing.T) {
			result := ParseCountryCode("+1 315-284-1580")

			assert.Equal(result, "US", "they should be equal")
		})

		t.Run("should parse country code correctly", func(t *testing.T) {
			result := ParseCountryCode("4566118311")

			assert.Equal(result, "DK", "they should be equal")
		})
	})

	t.Run("IsValid", func(t *testing.T) {
		t.Run("should validate phone number", func(t *testing.T) {
			result := IsValid("+1 315-284-1580")

			assert.Equal(result, true, "they should be equal")
		})

		t.Run("should validate phone number", func(t *testing.T) {
			result := IsValid("P+1 315-284-1580A")

			assert.Equal(result, false, "they should be equal")
		})
	})
}
