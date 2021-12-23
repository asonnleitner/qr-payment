package main

import (
	"fmt"
	"github.com/asonnleitner/qr-payment/iban"
)

func main() {
	accNum := "2171532/0800"

	IBAN, err := iban.ParseAccount(accNum, "CZ")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(IBAN)
}
