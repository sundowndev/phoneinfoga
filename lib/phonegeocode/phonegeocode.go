package phonegeocode

import (
	"errors"
	"strings"
	"sync"

	gotrie "github.com/tchap/go-patricia/patricia"
)

var ErrCountryNotFound = errors.New("could not identify country from phone number")

var (
	once = &sync.Once{}
	tree *trieGeocoder
)

type Geocoder interface {
	Country(number string) (string, error)
}

// init initialises a new thread-safe geocoder
func init() {
	once.Do(func() {
		tree = &trieGeocoder{
			data: initPrefixes(),
		}
	})
}

type trieGeocoder struct {
	data *gotrie.Trie
}

// Country tries to identify the country for a phone number - assuming it is provided in i18n format (+nn)
func Country(number string) (cc string, err error) {
	number = strings.TrimPrefix(number, "+")

	maxLen := -1
	_ = tree.data.VisitPrefixes(gotrie.Prefix(number), func(prefix gotrie.Prefix, item gotrie.Item) error {
		if len(prefix) > maxLen {
			cc = item.(string)
		}
		return nil
	})

	if len(cc) == 0 {
		err = ErrCountryNotFound
	}

	return
}
