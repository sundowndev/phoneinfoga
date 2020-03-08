# phone-iso3166

[![Go Report Card](https://goreportcard.com/badge/github.com/onlinecity/go-phone-iso3166)](https://goreportcard.com/report/github.com/onlinecity/go-phone-iso3166)
[![go-doc](https://godoc.org/github.com/onlinecity/go-phone-iso3166?status.svg)](https://godoc.org/github.com/onlinecity/go-phone-iso3166)

Maps an E.164 (international) phone number to the ISO-3166-1 alpha 2 (two letter) country code, associated with that number.

Also provides mapping for E.212 (mobile network codes, mcc+mnc) to the country.

It's based on it's python namesake [onlinecity/phone-iso3166](https://github.com/onlinecity/phone-iso3166) - but uses a [radix](https://github.com/hashicorp/go-immutable-radix) as it's datastructure.

### Usage

```go
import phoneiso3166 "github.com/onlinecity/go-phone-iso3166"
println(phoneiso3166.E164.Lookup(4566118311)) // prints: DK
```


### Performance

phone-iso3166 is reasonable fast

```
goos: darwin
goarch: amd64
pkg: github.com/onlinecity/go-phone-iso3166
BenchmarkE164Lookup-8           	20000000	       102 ns/op	      16 B/op	       1 allocs/op
BenchmarkE164LookupPrealloc-8   	20000000	        74.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkE164LookupBuffer-8     	20000000	       105 ns/op	      16 B/op	       1 allocs/op
BenchmarkE164LookupString-8     	20000000	        61.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkE164LookupBytes-8      	30000000	        53.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkE164LookupNoExist-8    	10000000	       151 ns/op	      16 B/op	       1 allocs/op
BenchmarkE212Lookup-8           	10000000	       163 ns/op	       3 B/op	       1 allocs/op
```

