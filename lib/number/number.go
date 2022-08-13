package number

import (
	"errors"
	"fmt"
	"github.com/nyaruka/phonenumbers"
	"github.com/sundowndev/phoneinfoga/v2/lib/phonegeocode"
)

// Number is a phone number
type Number struct {
	RawLocal      string
	Local         string
	E164          string
	International string
	CountryCode   int32
	Country       string
	Carrier       string
}

func NewNumber(number string) (res *Number, err error) {
	n := fmt.Sprintf("+%s", FormatNumber(number))

	country, err := phonegeocode.Country(number)
	if err != nil {
		return nil, errors.New("country code not recognized")
	}

	num, err := phonenumbers.Parse(n, "")
	if err != nil {
		return nil, err
	}

	res = &Number{
		RawLocal:      FormatNumber(phonenumbers.Format(num, phonenumbers.NATIONAL)),
		Local:         phonenumbers.Format(num, phonenumbers.NATIONAL),
		E164:          phonenumbers.Format(num, phonenumbers.E164),
		International: FormatNumber(phonenumbers.Format(num, phonenumbers.E164)),
		CountryCode:   num.GetCountryCode(),
		Country:       country,
		Carrier:       num.GetPreferredDomesticCarrierCode(),
	}

	return res, nil
}
