# Go module usage

You can easily use scanners in your own Golang script. You can find [Go documentation here](https://godoc.org/github.com/sundowndev/PhoneInfoga).

```go
package main

import (
        "fmt"
        
        "github.com/sundowndev/PhoneInfoga/pkg/scanners/numverify"
)

func main() {
        scan, err := numverify.NumverifyScan("<number>")

        if err != nil {
                log.Fatal(err)
        }

        fmt.Println(scan.CountryCode)
}
```