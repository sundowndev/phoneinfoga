# Go module usage

You can easily use scanners in your own Golang script. You can find [Go documentation here](https://pkg.go.dev/github.com/sundowndev/phoneinfoga/v2).

### Install the module

```
go get -v github.com/sundowndev/phoneinfoga/v2
```

### Usage example

```go
package main

import (
	"fmt"

	"github.com/sundowndev/phoneinfoga/v2/lib/number"
	"github.com/sundowndev/phoneinfoga/v2/lib/remote"
)

func main() {
	n, _ := number.NewNumber("...")
	s := remote.NewGoogleSearchScanner()

	res, _ := s.Scan(n)

	links := res.(remote.GoogleSearchResponse)
	for _, link := range links.Individuals {
		fmt.Println(link.URL) // Google search link to scan
	}
}
```
