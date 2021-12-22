package main

import (
	"fmt"
	"github.com/asonnleitner/qr-payment/iban"
)

func main() {
	IBAN := iban.ParseAccount("2171532/0800", "CZ")

	fmt.Println(IBAN)
}
