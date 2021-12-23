package number

import (
	"regexp"
	"strconv"

	phoneiso3166 "github.com/onlinecity/go-phone-iso3166"
)

// FormatNumber formats a phone number to remove
// unnecessary chars and avoid dealing with unwanted input.
func FormatNumber(n string) string {
	re := regexp.MustCompile(`[_\W]+`)
	number := re.ReplaceAllString(n, "")

	return number
}

// ParseCountryCode parses a phone number and returns ISO country code.
// This is required in order to use the phonenumbers library.
func ParseCountryCode(n string) string {
	var number uint64
	number, _ = strconv.ParseUint(FormatNumber(n), 10, 64)

	return phoneiso3166.E164.Lookup(number)
}

// IsValid indicate if a phone number has a valid format.
func IsValid(number string) bool {
	number = FormatNumber(number)

	re := regexp.MustCompile("^[0-9]+$")

	return len(re.FindString(number)) != 0
}
