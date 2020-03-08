//+build debug

package oncer

import (
	"fmt"
	"time"
)

func log(name string, fn func()) func() {
	return func() {
		start := time.Now()
		if len(name) > 80 {
			name = name[(len(name) - 80):]
		}
		defer fmt.Println(name, time.Now().Sub(start))
		fn()
	}
}
