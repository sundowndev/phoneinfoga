# Phone Geocode

**This package was copied from [github.com/davegardnerisme/phonegeocode](https://github.com/davegardnerisme/phonegeocode) then modified.**

----

Internationalised phone number geocoding by country for Go.

I built this because I needed to turn phone numbers into countries, and that's
really _all_ I needed. If [libphonenumber](https://code.google.com/p/libphonenumber/)
existed for Go, I would probably just use that. AFAIK it doesn't.

This is based on work in [github.com/relops/prefixes](https://github.com/relops/prefixes),
however it has a different implementation, using a [Trie](http://en.wikipedia.org/wiki/Trie)
data structure - specifically [github.com/tchap/go-patricia](https://github.com/tchap/go-patricia).

The way it works is that we have a list of prefixes that identify a country, and
we simply match the _most specific prefix_ to find the country code. This deals
with Canada where the country code is `+1` and shared with US.

All the data lives in a CSV and can be used via codegen to create our Trie.

```
go run data/codegen.go > ./data.go && go fmt
```

## Usage

```
// cc = GB, err = nil
cc, err := phonegeocode.Country("+447999111222")

// cc = "", err = phonegeocode.ErrCountryNotFound
cc, err := phonegeocode.Country("+999999999998")
```

## License

The MIT License (MIT)

Copyright (c) 2016 Dave Gardner

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
