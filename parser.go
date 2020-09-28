package main

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

// groupMap gives a map with as keys the group names and as values the matches
func groupMap(r *regexp.Regexp, s string) map[string]string {
	result := make(map[string]string)
	keys := r.SubexpNames()
	values := r.FindStringSubmatch(s)
	for i, value := range values {
		if i != 0 {
			result[keys[i]] = value
		}
	}
	return result
}

// trimValues trims the spaces of all the values in the map
func trimValues(m map[string]string) {
	for key, value := range m {
		m[key] = strings.TrimSpace(value)
	}
}

var initialRecordRegex = regexp.MustCompile(`^\d{1}` +
	`.{4}` +
	`(?P<creation_date>\d{6})` +
	`(?P<bank_identification_number>\d{3})` +
	`.{2}` +
	`(?P<duplicate>[D\s]{1})` +
	`.{7}` +
	`(?P<reference>.{10})` +
	`(?P<addressee>.{26})` +
	`(?P<bic>.{11})` +
	`0?(?P<account_holder_reference>\d{10})` +
	`.{1}` +
	`(?P<free>.{5})` +
	`(?P<transaction_reference>.{16})` +
	`(?P<related_reference>.{16})` +
	`.{7}` +
	`(?P<version_code>\d{1})$`)

var oldBalanceRecordRegex = regexp.MustCompile(`^\d{1}` +
	`(?P<account_structure>\d{1})` +
	`(?P<serial_number>\d{3})` +
	`(?P<account_number>.{37})` +
	`(?P<balance_sign>[01]{1})` +
	`(?P<old_balance>\d{15})` +
	`(?P<balance_date>\d{6})` +
	`(?P<account_holder_name>.{26})` +
	`(?P<account_description>.{35})` +
	`(?P<bank_statement_serial_number>\d{3})$`)

// Record represents a generic line in a CODA file
type Record interface {
	Parse(string) error
}

// InitialRecord represents the first line of the CODA file
type InitialRecord struct {
	CreationDate             time.Time
	BankIdentificationNumber int
	IsDuplicate              bool
	Reference                string
	Addressee                string
	BIC                      string
	AccountHolderReference   int
	Free                     string
	TransactionReference     string
	VersionCode              int
}

// OldBalanceRecord represents the old balance at the start of the CODA file
type OldBalanceRecord struct {
	AccountStructure          int
	SerialNumber              int
	AccountNumber             string
	BalanceSign               bool // True means debit / false credit
	OldBalance                int
	BalanceDate               time.Time
	AccountHolderName         string
	AccountDescription        string
	BankStatementSerialNumber int
}

// Parse reads data from string s into an initial record
func (r *InitialRecord) Parse(s string) (err error) {
	m := groupMap(initialRecordRegex, s)
	trimValues(m)

	r.CreationDate, err = time.Parse("020106", m["creation_date"])
	if err != nil {
		return err
	}
	if m["bank_identification_number"] != "" {
		r.BankIdentificationNumber, err = strconv.Atoi(m["bank_identification_number"])
		if err != nil {
			return err
		}
	}
	r.IsDuplicate = m["duplicate"] == "D"
	r.Reference = m["reference"]
	r.Addressee = m["addressee"]
	r.BIC = m["bic"]
	if m["account_holder_reference"] != "" {
		r.AccountHolderReference, err = strconv.Atoi(m["account_holder_reference"])
		if err != nil {
			return err
		}
	}
	r.Free = m["free"]
	r.TransactionReference = m["transaction_reference"]
	if m["version_code"] != "" {
		r.VersionCode, err = strconv.Atoi(m["version_code"])
		if err != nil {
			return err
		}
	}
	return nil
}

// Parse reads data from string s into an old balance record
func (r *OldBalanceRecord) Parse(s string) (err error) {
	m := groupMap(oldBalanceRecordRegex, s)
	trimValues(m)

	if m["account_structure"] != "" {
		r.AccountStructure, err = strconv.Atoi(m["account_structure"])
		if err != nil {
			return err
		}
	}
	if m["serial_number"] != "" {
		r.SerialNumber, err = strconv.Atoi(m["serial_number"])
		if err != nil {
			return err
		}
	}
	r.AccountNumber = m["account_number"]
	r.BalanceSign = m["balance_sign"] == "1"
	if m["old_balance"] != "" {
		r.OldBalance, err = strconv.Atoi(m["old_balance"])
		if err != nil {
			return err
		}
	}
	r.BalanceDate, err = time.Parse("020106", m["balance_date"])
	if err != nil {
		return err
	}
	r.AccountHolderName = m["account_holder_name"]
	r.AccountDescription = m["account_description"]
	if m["bank_statement_serial_number"] != "" {
		r.BankStatementSerialNumber, err = strconv.Atoi(m["bank_statement_serial_number"])
		if err != nil {
			return err
		}
	}
	return nil
}

// Parse takes a line of CODA and parses it
func Parse(line string) (r Record, err error) {
	if strings.HasPrefix(line, "0") {
		r = &InitialRecord{}
	} else if strings.HasPrefix(line, "1") {
		r = &OldBalanceRecord{}
	} else {
		return nil, nil
	}

	err = r.Parse(line)
	return r, err
}
