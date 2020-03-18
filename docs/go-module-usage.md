# Go module usage

You can easily use scanners in your own Golang script. You can find [Go documentation here](https://godoc.org/github.com/sundowndev/PhoneInfoga).

```go
package main

import (
	"fmt"
	"log"

	phoneinfoga "github.com/sundowndev/phoneinfoga/pkg/scanners"
)

func main() {
	number, err := phoneinfoga.LocalScan(number)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(number.E164)
}
```
