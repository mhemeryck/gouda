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

var initialRecordRegex = regexp.MustCompile(`^0` +
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

var oldBalanceRecordRegex = regexp.MustCompile(`^1` +
	`(?P<account_structure>\d{1})` +
	`(?P<serial_number>\d{3})` +
	`(?P<account_number>.{37})` +
	`(?P<balance_sign>[01]{1})` +
	`(?P<old_balance>\d{15})` +
	`(?P<balance_date>\d{6})` +
	`(?P<account_holder_name>.{26})` +
	`(?P<account_description>.{35})` +
	`(?P<bank_statement_serial_number>\d{3})$`)

var transactionRecordRegex = regexp.MustCompile(`^21` +
	`(?P<serial_number>\d{4})` +
	`(?P<detail_number>\d{4})` +
	`(?P<bank_reference_number>.{21})` +
	`(?P<balance_sign>\d{1})` +
	`(?P<balance>\d{15})` +
	`(?P<balance_date>\d{6})` +
	`(?P<transaction_code>\d{8})` +
	`(?P<reference_type>\d{1})` +
	`(?P<reference>.{53})` +
	`(?P<booking_date>\d{6})` +
	`(?P<bank_statement_serial_number>\d{3})` +
	`(?P<globalisation_code>\d{1})` +
	`(?P<transaction_sequence>[01]{1})` +
	`\s{1}` +
	`(?P<information_sequence>[01]{1})$`)

var transactionPurposeRecordRegex = regexp.MustCompile(`^22` +
	`(?P<serial_number>\d{4})` +
	`(?P<detail_number>\d{4})` +
	`(?P<bank_statement>.{53})` +
	`(?P<client_reference>.{35})` +
	`(?P<bic>.{11})` +
	`\s{3}` +
	`(?P<transaction_type>[\s12345]{1})` +
	`(?P<reason_return_code>.{4})` +
	`(?P<purpose_category>.{4})` +
	`(?P<purpose>.{4})` +
	`(?P<transaction_sequence>[01]{1})` +
	`\s{1}` +
	`(?P<information_sequence>[01]{1})$`)

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

// TransactionRecord represents a single transaction in a CODA file
type TransactionRecord struct {
	SerialNumber              int
	DetailNumber              int
	BankReferenceNumber       string
	BalanceSign               bool // True means debit / false credit
	Balance                   int
	BalanceDate               time.Time
	TransactionCode           int
	ReferenceType             int
	Reference                 string
	BookingDate               time.Time
	BankStatementSerialNumber int
	GlobalisationCode         int
	TransactionSequence       bool
	InformationSequence       bool
}

// TransactionPurposeRecord holds extra information related to the transaction record
type TransactionPurposeRecord struct {
	SerialNumber        int
	DetailNumber        int
	BankStatement       string
	ClientReference     string
	BIC                 string
	TransactionType     int
	ReasonReturnCode    string
	PurposeCategory     string
	Purpose             string
	TransactionSequence bool
	InformationSequence bool
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
	return err
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
	return err
}

// Parse reads data from string s into a transaction record
func (r *TransactionRecord) Parse(s string) (err error) {
	m := groupMap(transactionRecordRegex, s)
	trimValues(m)

	if m["serial_number"] != "" {
		r.SerialNumber, err = strconv.Atoi(m["serial_number"])
		if err != nil {
			return err
		}
	}
	if m["detail_number"] != "" {
		r.DetailNumber, err = strconv.Atoi(m["detail_number"])
		if err != nil {
			return err
		}
	}
	r.BankReferenceNumber = m["bank_reference_number"]
	r.BalanceSign = m["balance_sign"] == "1"
	if m["balance"] != "" {
		r.Balance, err = strconv.Atoi(m["balance"])
		if err != nil {
			return err
		}
	}
	if m["balance_date"] != "000000" {
		r.BalanceDate, err = time.Parse("020106", m["balance_date"])
		if err != nil {
			return err
		}
	}
	if m["transaction_code"] != "" {
		r.TransactionCode, err = strconv.Atoi(m["transaction_code"])
		if err != nil {
			return err
		}
	}
	if m["reference_type"] != "" {
		r.ReferenceType, err = strconv.Atoi(m["reference_type"])
		if err != nil {
			return err
		}
	}
	r.Reference = m["reference"]
	r.BookingDate, err = time.Parse("020106", m["booking_date"])
	if err != nil {
		return err
	}
	if m["bank_statement_serial_number"] != "" {
		r.BankStatementSerialNumber, err = strconv.Atoi(m["bank_statement_serial_number"])
		if err != nil {
			return err
		}
	}
	if m["globalisation_code"] != "" {
		r.GlobalisationCode, err = strconv.Atoi(m["globalisation_code"])
		if err != nil {
			return err
		}
	}
	r.TransactionSequence = m["transaction_sequence"] == "1"
	r.InformationSequence = m["information_sequence"] == "1"

	return err
}

// Parse will wrap a transaction purpose record line
func (r *TransactionPurposeRecord) Parse(s string) (err error) {
	m := groupMap(transactionPurposeRecordRegex, s)
	trimValues(m)

	if m["serial_number"] != "" {
		r.SerialNumber, err = strconv.Atoi(m["serial_number"])
		if err != nil {
			return err
		}
	}
	if m["detail_number"] != "" {
		r.DetailNumber, err = strconv.Atoi(m["detail_number"])
		if err != nil {
			return err
		}
	}
	r.BankStatement = m["bank_statement"]
	r.ClientReference = m["client_reference"]
	r.BIC = m["bic"]
	if m["transaction_type"] != "" {
		r.TransactionType, err = strconv.Atoi(m["transaction_type"])
		if err != nil {
			return err
		}
	}
	r.ReasonReturnCode = m["reason_return_code"]
	r.PurposeCategory = m["purpose_category"]
	r.Purpose = m["Purpose"]
	r.TransactionSequence = m["transaction_sequence"] == "1"
	r.InformationSequence = m["information_sequence"] == "1"

	return err
}

// Parse takes a line of CODA and parses it
func Parse(line string) (r Record, err error) {
	if strings.HasPrefix(line, "0") {
		r = &InitialRecord{}
	} else if strings.HasPrefix(line, "1") {
		r = &OldBalanceRecord{}
	} else if strings.HasPrefix(line, "21") {
		r = &TransactionRecord{}
	} else if strings.HasPrefix(line, "22") {
		r = &TransactionPurposeRecord{}
	} else {
		return nil, nil
	}

	err = r.Parse(line)
	return r, err
}
