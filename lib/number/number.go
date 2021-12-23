package number

import (
	"github.com/nyaruka/phonenumbers"
)

// Number is a phone number
type Number struct {
	RawLocal      string `json:"rawLocal"`
	Local         string `json:"local"`
	E164          string `json:"E164"`
	International string `json:"international"`
	CountryCode   int32  `json:"countryCode"`
	Country       string `json:"country"`
	Carrier       string `json:"carrier"`
}

func NewNumber(number string) (res *Number, err error) {
	n := "+" + FormatNumber(number)
	country := ParseCountryCode(n)

	num, err := phonenumbers.Parse(n, country)
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
