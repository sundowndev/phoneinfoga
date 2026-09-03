package test

import (
	"log"

	"github.com/sundowndev/phoneinfoga/v2/lib/number"
)

func NewFakeUSNumber() *number.Number {
	n, err := number.NewNumber("+1.4152229670")
	if err != nil {
		log.Fatalf("NewFakeUSNumber: %s", err)
	}
	return n
}

func NewFakeFRNumber() *number.Number {
	n, err := number.NewNumber("+33.679358133")
	if err != nil {
		log.Fatalf("NewFakeFRNumber: %s", err)
	}
	return n
}
