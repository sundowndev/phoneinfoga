//+build !debug

package oncer

func log(name string, fn func()) func() {
	return fn
}
