package number

import (
	"github.com/nyaruka/phonenumbers"
)

// Number is a phone number
type Number struct {
	Valid         bool
	RawLocal      string
	Local         string
	E164          string
	International string
	CountryCode   int32
	Country       string
	Carrier       string
}

func NewNumber(number string) (res *Number, err error) {
	n := "+" + FormatNumber(number)
	country := ParseCountryCode(n)

	num, err := phonenumbers.Parse(n, country)
	if err != nil {
		return nil, err
	}

	res = &Number{
		Valid:         phonenumbers.IsValidNumber(num),
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
