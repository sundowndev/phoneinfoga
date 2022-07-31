package number

import (
	_ "embed"
	"regexp"
)

// FormatNumber formats a phone number to remove
// unnecessary chars and avoid dealing with unwanted input.
func FormatNumber(n string) string {
	re := regexp.MustCompile(`[_\W]+`)
	return re.ReplaceAllString(n, "")
}

// IsValid indicate if a phone number has a valid format.
func IsValid(number string) bool {
	number = FormatNumber(number)

	re := regexp.MustCompile("^[0-9]+$")

	return len(re.FindString(number)) != 0
}
