package utils

import (
	"testing"

	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
	"gopkg.in/sundowndev/phoneinfoga.v2/pkg/utils/mocks"
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

	t.Run("Logger", func(t *testing.T) {
		t.Run("Infoln", func(t *testing.T) {
			mLogger := new(mocks.Color)

			mLogger.On("Println", "[i]", "test").Once().Return(0, nil)

			log := &Logger{
				NewColor: func(value ...color.Attribute) colorLogger {
					assert.Equal([]color.Attribute([]color.Attribute{36}), value, "they should be equal")

					return mLogger
				},
			}

			log.Infoln("test")

			mLogger.AssertExpectations(t)
		})

		t.Run("Warnln", func(t *testing.T) {
			mLogger := new(mocks.Color)

			mLogger.On("Println", "[*]", "test").Once().Return(0, nil)

			log := &Logger{
				NewColor: func(value ...color.Attribute) colorLogger {
					assert.Equal([]color.Attribute([]color.Attribute{33}), value, "they should be equal")

					return mLogger
				},
			}

			log.Warnln("test")

			mLogger.AssertExpectations(t)
		})

		t.Run("Errorln", func(t *testing.T) {
			mLogger := new(mocks.Color)

			mLogger.On("Println", "[!]", "test").Once().Return(0, nil)

			log := &Logger{
				NewColor: func(value ...color.Attribute) colorLogger {
					assert.Equal([]color.Attribute([]color.Attribute{31}), value, "they should be equal")

					return mLogger
				},
			}

			log.Errorln("test")

			mLogger.AssertExpectations(t)
		})

		t.Run("Successln", func(t *testing.T) {
			mLogger := new(mocks.Color)

			mLogger.On("Println", "[+]", "test").Once().Return(0, nil)

			log := &Logger{
				NewColor: func(value ...color.Attribute) colorLogger {
					assert.Equal([]color.Attribute([]color.Attribute{32}), value, "they should be equal")

					return mLogger
				},
			}

			log.Successln("test")

			mLogger.AssertExpectations(t)
		})

		t.Run("Successf", func(t *testing.T) {
			var ColorNumberOfCalls int

			mLogger := new(mocks.Color)

			mLogger.On("Printf", "[+] %s", "test").Once().Return(0, nil)
			mLogger.On("Printf", "\n").Once().Return(0, nil)

			log := &Logger{
				NewColor: func(value ...color.Attribute) colorLogger {
					if ColorNumberOfCalls == 0 {
						assert.Equal([]color.Attribute([]color.Attribute{32}), value, "they should be equal")
						ColorNumberOfCalls++
					}

					return mLogger
				},
			}

			log.Successf("%s", "test")

			mLogger.AssertNumberOfCalls(t, "Printf", 2)
			mLogger.AssertExpectations(t)
		})
	})
}
