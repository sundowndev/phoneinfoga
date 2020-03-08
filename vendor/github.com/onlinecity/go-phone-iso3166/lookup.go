package phoneiso3166

import (
	"strconv"

	iradix "github.com/hashicorp/go-immutable-radix"
)

//go:generate python3 gen/e164.py e164.go
//go:generate python3 gen/e212.py e212.go

type CountryRadix struct {
	*iradix.Tree
}
type OperatorRadix struct {
	*iradix.Tree
}

var E164 *CountryRadix
var E212 *OperatorRadix

func init() {
	E164 = &CountryRadix{
		getE164(),
	}
	E212 = &OperatorRadix{
		getE212(),
	}
}

func (r *CountryRadix) Lookup(msisdn uint64) string {
	return r.LookupByteString([]byte(strconv.FormatUint(msisdn, 10)))
}

func (r *CountryRadix) LookupString(msisdn string) string {
	return r.LookupByteString([]byte(msisdn))
}

func (r *CountryRadix) LookupByteString(msisdn []byte) string {
	_, country, _ := r.Root().LongestPrefix(msisdn)
	if country == nil {
		return ""
	}
	return country.(string)
}

func (r *OperatorRadix) Lookup(mcc, mnc uint16) string {
	c := strconv.FormatUint(uint64(mcc), 10)
	n := strconv.FormatUint(uint64(mnc), 10)
	_, country, _ := r.Root().LongestPrefix([]byte(c + n))
	if country == nil {
		return ""
	}
	return country.(string)
}

func NetworkName(mcc, mnc uint16) string {
	if op, ok := OperatorMap[mcc][mnc]; ok {
		return op.Name
	}
	return ""
}
