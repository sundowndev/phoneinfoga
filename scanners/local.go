package scanners

import (
	"github.com/nyaruka/phonenumbers"
	"gopkg.in/sundowndev/phoneinfoga.v2/utils"
)

// LocalScan performs a local scan of a phone
// number using phonenumbers library
func LocalScan(number string) (res *Number, err error) {
	n := "+" + utils.FormatNumber(number)
	country := utils.ParseCountryCode(n)

	num, err := phonenumbers.Parse(n, country)
	if err != nil {
		return nil, err
	}

	res = &Number{
		RawLocal:      utils.FormatNumber(phonenumbers.Format(num, phonenumbers.NATIONAL)),
		Local:         phonenumbers.Format(num, phonenumbers.NATIONAL),
		E164:          phonenumbers.Format(num, phonenumbers.E164),
		International: utils.FormatNumber(phonenumbers.Format(num, phonenumbers.E164)),
		CountryCode:   num.GetCountryCode(),
		Country:       country,
		Carrier:       num.GetPreferredDomesticCarrierCode(),
	}

	return res, nil
}
