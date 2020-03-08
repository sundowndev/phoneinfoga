package scanners

import (
	"github.com/nyaruka/phonenumbers"
	"github.com/sundowndev/phoneinfoga/pkg/utils"
)

// LocalScan performs a local scan of a phone
// number using phonenumbers library
func LocalScan(number string) (res *Number, err error) {
	country := utils.ParseCountryCode(number)

	num, err := phonenumbers.Parse(number, country)

	if err != nil {
		return nil, err
	}

	res = &Number{
		Local:         phonenumbers.Format(num, phonenumbers.NATIONAL),
		E164:          phonenumbers.Format(num, phonenumbers.E164),
		International: utils.FormatNumber(phonenumbers.Format(num, phonenumbers.E164)),
		CountryCode:   num.GetCountryCode(),
		Country:       country,
		Carrier:       num.GetPreferredDomesticCarrierCode(),
	}

	return res, nil
}
