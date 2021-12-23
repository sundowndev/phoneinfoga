package scanners

import (
	"github.com/nyaruka/phonenumbers"
	libnum "github.com/sundowndev/phoneinfoga/v2/lib/number"
)

// LocalScan performs a local scan of a phone
// number using phonenumbers library
func LocalScan(number string) (res *Number, err error) {
	n := "+" + libnum.FormatNumber(number)
	country := libnum.ParseCountryCode(n)

	num, err := phonenumbers.Parse(n, country)
	if err != nil {
		return nil, err
	}

	res = &Number{
		RawLocal:      libnum.FormatNumber(phonenumbers.Format(num, phonenumbers.NATIONAL)),
		Local:         phonenumbers.Format(num, phonenumbers.NATIONAL),
		E164:          phonenumbers.Format(num, phonenumbers.E164),
		International: libnum.FormatNumber(phonenumbers.Format(num, phonenumbers.E164)),
		CountryCode:   num.GetCountryCode(),
		Country:       country,
		Carrier:       num.GetPreferredDomesticCarrierCode(),
	}

	return res, nil
}
