# Go module usage

You can easily use scanners in your own Golang script. You can find [Go documentation here](https://pkg.go.dev/gopkg.in/sundowndev/phoneinfoga.v2).

### Install the module

```
go get -v gopkg.in/sundowndev/phoneinfoga.v2
```

### Usage example

```go
package main

import (
	"fmt"
	"log"

	phoneinfoga "gopkg.in/sundowndev/phoneinfoga.v2/pkg/scanners"
)

func main() {
	number, err := phoneinfoga.LocalScan("<number>")

	if err != nil {
		log.Fatal(err)
	}

	links := phoneinfoga.GoogleSearchScan(number)

	for _, link := range links.Individuals {
		fmt.Println(link.URL) // Google search link to scan
	}
}
```
