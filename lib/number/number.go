package number

import (
	"bytes"
	"fmt"
	"github.com/nyaruka/phonenumbers"
	"strings"
	"text/template"
)

const templateDigitPlaceholder = "x"

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
	CustomFormats []string
}

type FormatTemplateData struct {
	CountryCode string
	Country     string
}

func NewNumber(number string, formats ...string) (res *Number, err error) {
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
		CustomFormats: []string{},
	}

	for _, format := range formats {
		res.CustomFormats = append(res.CustomFormats, convertFormatTemplate(res, format))
	}

	return res, nil
}

func convertFormatTemplate(n *Number, f string) string {
	countryCodeStr := fmt.Sprintf("%d", n.CountryCode)

	var out []string
	splitNumber := strings.Split(strings.Replace(n.International, countryCodeStr, "", 1), "")
	splitFormat := strings.Split(f, "")
	for _, j := range splitFormat {
		if strings.ToLower(j) == templateDigitPlaceholder && len(splitNumber) > 0 {
			j = splitNumber[0]
			splitNumber = splitNumber[1:]
		}
		out = append(out, j)
	}

	t := template.Must(template.New("custom-format").Parse(strings.Join(out, "")))

	var tpl bytes.Buffer
	err := t.Execute(&tpl, FormatTemplateData{
		CountryCode: countryCodeStr,
		Country:     n.Country,
	})
	if err != nil {
		return strings.Join(out, "")
	}
	return tpl.String()
}
