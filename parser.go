package main

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

// groupMap gives a map with as keys the group names and as values the matches
func groupMap(r *regexp.Regexp, s string) map[string]*string {
	result := make(map[string]*string)
	keys := r.SubexpNames()
	values := r.FindStringSubmatch(s)
	for i := range values {
		if i != 0 {
			result[keys[i]] = &values[i]
		}
	}
	return result
}

// parseInt is a helper method to parse an integer from a given string
func parseInt(s string) (*int, error) {
	s = strings.TrimSpace(s)
	// Nothing found, just point to nothing
	if s == "" {
		return nil, nil
	}
	v, err := strconv.Atoi(s)
	return &v, err
}

// trimValues trims the spaces of all the values in the map
func trimValues(m map[string]*string) {
	for i := range m {
		*m[i] = strings.TrimSpace(*m[i])
	}
}

var initialRecordRegex = regexp.MustCompile(`^[\d\s]{1}` +
	`.{4}` +
	`(?P<creation_date>\d{6})` +
	`(?P<bank_identification_number>[\d\s]{3})` +
	`.{2}` +
	`(?P<duplicate>.{1})` +
	`.{7}` +
	`(?P<reference>.{10})` +
	`(?P<addressee>.{26})` +
	`(?P<bic>.{11})` +
	`0?(?P<account_holder_reference>[\d\s]{10})` +
	`.{1}` +
	`(?P<free>.{5})` +
	`(?P<transaction_reference>.{16})` +
	`(?P<related_reference>.{16})` +
	`.{7}` +
	`(?P<version_code>[\d\s]{1})$`)

var oldBalanceRecordRegex = regexp.MustCompile(`^[\d\s]{1}` +
	`(?P<account_structure>[\d\s]{1})` +
	`(?P<serial_number>[\d\s]{3})` +
	`(?P<account_number>.{37})` +
	`(?P<balance_sign>[\d\s]{1})` +
	`(?P<old_balance>\d{15})` +
	`(?P<balance_date>\d{6})` +
	`(?P<account_holder_name>.{26})` +
	`(?P<account_description>.{35})` +
	`(?P<bank_statement_serial_number>[\d\s]{3})$`)

// InitialRecord represents the first line of the CODA file
type InitialRecord struct {
	CreationDate             *time.Time
	BankIdentificationNumber *int
	IsDuplicate              *bool
	Reference                *string
	Addressee                *string
	BIC                      *string
	AccountHolderReference   *int
	Free                     *string
	TransactionReference     *string
	VersionCode              *int
}

// OldBalanceRecord represents the old balance at the start of the CODA file
type OldBalanceRecord struct {
	AccountStructure          *int
	SerialNumber              *int
	AccountNumber             *string
	BalanceSign               *bool // True means debit / false credit
	OldBalance                *int
	BalanceDate               *time.Time
	AccountHolderName         *string
	AccountDescription        *string
	BankStatementSerialNumber *int
}

// ParseInitialRecord creates an InitialRecord from a given string
func ParseInitialRecord(s string) (r InitialRecord, err error) {
	m := groupMap(initialRecordRegex, s)
	trimValues(m)

	v, err := time.Parse("020106", *m["creation_date"])
	if err != nil {
		return r, err
	}
	r.CreationDate = &v
	r.BankIdentificationNumber, err = parseInt(*m["bank_identification_number"])
	if err != nil {
		return r, err
	}
	isDuplicate := *m["duplicate"] == "D"
	r.IsDuplicate = &isDuplicate
	r.Reference = m["reference"]
	r.Addressee = m["addressee"]
	r.BIC = m["bic"]
	r.AccountHolderReference, err = parseInt(*m["account_holder_reference"])
	if err != nil {
		return r, err
	}
	r.Free = m["free"]
	r.TransactionReference = m["transaction_reference"]
	r.VersionCode, err = parseInt(*m["version_code"])
	if err != nil {
		return r, err
	}
	return r, nil
}

// ParseOldBalanceRecord creates an OldBalanceRecord from a string
func ParseOldBalanceRecord(s string) (r OldBalanceRecord, err error) {
	m := groupMap(oldBalanceRecordRegex, s)
	trimValues(m)

	r.AccountStructure, err = parseInt(*m["account_structure"])
	if err != nil {
		return r, err
	}
	r.SerialNumber, err = parseInt(*m["serial_number"])
	if err != nil {
		return r, err
	}
	r.AccountNumber = m["account_number"]
	balanceSign := *m["balance_sign"] == "1"
	r.BalanceSign = &balanceSign
	r.OldBalance, err = parseInt(*m["old_balance"])
	if err != nil {
		return r, err
	}
	v, err := time.Parse("020106", *m["balance_date"])
	if err != nil {
		return r, err
	}
	r.BalanceDate = &v
	r.AccountHolderName = m["account_holder_name"]
	r.AccountDescription = m["account_description"]
	r.BankStatementSerialNumber, err = parseInt(*m["bank_statement_serial_number"])
	return r, nil
}
