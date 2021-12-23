package iban

import (
	"fmt"
	"github.com/asonnleitner/qr-payment/strutils"
	"github.com/asonnleitner/qr-payment/utils"
	"regexp"
	"strconv"
	"strings"
)

const (
	errInvalidAccountNumber    = "invalid account number"
	errCountryCodeNotSupported = "country code not supported"
)

var (
	czBankAccRegExp = regexp.MustCompile(`^([0-9]{2,6})-([0-9]{2,10})/([0-9]{4})$|^([0-9]{2,10})/([0-9]{4})$`)
	lettersPattern  = regexp.MustCompile(`[A-Z]`)
)

func ParseAccount(account, country string) (string, error) {
	var pattern *regexp.Regexp
	country = strings.ToUpper(country)
	account = strings.Replace(account, " ", "", -1)

	switch country {
	case "CZ":
		pattern = czBankAccRegExp
	default:
		panic(errCountryCodeNotSupported)
	}

	prefix, number, bankCode := splitAccount(account, pattern)

	return convertToIBAN(country, prefix, number, bankCode)
}

func convertToIBAN(country, prefix, number, bankCode string) (string, error) {
	var bban string

	switch country {
	case "CZ":
		bban = strutils.Concat(4+6+10, bankCode, prefix, number)
	default:
		return "", fmt.Errorf(errCountryCodeNotSupported)
	}

	checksum := calculateChecksum(country, bban)

	return strutils.Concat(2+2+10, country, checksum, bban), nil
}

func calculateChecksum(country, bban string) string {
	var sb strings.Builder

	sb.WriteString(bban)
	sb.WriteString(country)
	sb.WriteString("00")

	checksum := 98 - utils.Modulo(
		replaceLetters(sb.String()), 97,
	)

	return strconv.FormatInt(checksum, 10)
}

func splitAccount(account string, pattern *regexp.Regexp) (prefix, number, bankCode string) {
	account = strings.TrimSpace(account)

	if !strings.Contains(account, "-") {
		account = "00-" + account
	}

	if len(account) < 10 || len(account) > 22 {
		panic(errInvalidAccountNumber)
	}

	parts := pattern.FindStringSubmatch(account)[1:]

	prefix = utils.Padding(
		strings.Replace(parts[0], "-", "", -1), 6, "0", utils.PaddingLeft,
	)
	number = utils.Padding(
		strings.Replace(parts[1], "/", "", -1), 10, "0", utils.PaddingLeft,
	)
	bankCode = parts[2]
	return
}

func replaceLetters(iban string) string {
	letters := lettersPattern.FindAllString(iban, -1)

	for i := 0; i < len(letters); i++ {
		iban = strings.Replace(iban, letters[i], fmt.Sprintf("%d", letters[i][0]-'A'+10), -1)
	}

	return iban
}
