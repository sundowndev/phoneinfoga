package test

import "github.com/sundowndev/phoneinfoga/v2/lib/number"

func NewFakeUSNumber() *number.Number {
	n, _ := number.NewNumber("+1.4152229670")
	return n
}
