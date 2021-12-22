package iban

import (
	"fmt"
	"github.com/asonnleitner/qr-payment/utils"
	"regexp"
	"strconv"
	"strings"
)

const (
	errInvalidAccountNumber    = "Invalid account number"
	errCountryCodeNotSupported = "Country code not supported"
)

var (
	czechBankAccountPattern = regexp.MustCompile(`^([0-9]{2,6}-)?([0-9]{2,10}/)([0-9]{4})$`)
	lettersPattern          = regexp.MustCompile(`[A-Z]`)
)

func ParseAccount(account, country string) string {
	var pattern *regexp.Regexp
	country = strings.ToUpper(country)
	account = strings.Replace(account, " ", "", -1)

	switch country {
	case "CZ":
		pattern = czechBankAccountPattern
	default:
		panic(errCountryCodeNotSupported)
	}

	prefix, number, bankCode := splitAccount(account, pattern)

	return convertToIBAN(country, prefix, number, bankCode)
}

func convertToIBAN(country, prefix, number, bankCode string) string {
	var sb strings.Builder

	switch country {
	case "CZ":
		sb.WriteString(bankCode)
		sb.WriteString(prefix)
		sb.WriteString(number)
	default:
		panic(errCountryCodeNotSupported)
	}

	bban := sb.String()

	checksum := calculateChecksum(country, bban)

	sb.Reset()
	sb.WriteString(country)
	sb.WriteString(checksum)
	sb.WriteString(bban)

	return sb.String()
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

	parts := pattern.FindStringSubmatch(account)

	prefix = utils.Padding(
		strings.Replace(parts[1], "-", "", -1), 6, "0", utils.PaddingLeft,
	)
	number = utils.Padding(
		strings.Replace(parts[2], "/", "", -1), 10, "0", utils.PaddingLeft,
	)
	bankCode = parts[3]
	return
}

func replaceLetters(iban string) string {
	letters := lettersPattern.FindAllString(iban, -1)

	for i := 0; i < len(letters); i++ {
		iban = strings.Replace(iban, letters[i], fmt.Sprintf("%d", letters[i][0]-'A'+10), -1)
	}

	return iban
}
