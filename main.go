package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"
)

const sample = `0000016032064505                  J. VRANCKEN BV            JVBABE22   00860332194 00005                                       2`

var initialRecordRegex = regexp.MustCompile(`^(?P<identification>[\d\s]{1})` +
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
	`(?P<version_code>[\d\s]{1})$`,
)

// InitialRecord represents the first line of the CODA file
type InitialRecord struct {
	Identification           int
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

// ParseInitialRecord creates an InitialRecord from a given string
func ParseInitialRecord(s string) (r InitialRecord, err error) {
	m := groupMap(initialRecordRegex, s)

	r.Identification, err = strconv.Atoi(m["identification"])
	if err != nil {
		return r, err
	}
	r.CreationDate, err = time.Parse("020106", m["creation_date"])
	if err != nil {
		return r, err
	}
	r.BankIdentificationNumber, err = strconv.Atoi(m["bank_identification_number"])
	if err != nil {
		return r, err
	}
	r.IsDuplicate = m["duplicate"] == "D"
	r.Reference = m["reference"]
	r.Addressee = m["addressee"]
	r.BIC = m["bic"]
	r.AccountHolderReference, err = strconv.Atoi(m["account_holder_reference"])
	if err != nil {
		return r, err
	}

	r.Free = m["free"]
	r.TransactionReference = m["transaction_reference"]
	r.VersionCode, err = strconv.Atoi(m["version_code"])
	if err != nil {
		return r, err
	}

	return r, nil
}

// groupMap gives a map with as keys the group names and as values the matches
func groupMap(r *regexp.Regexp, s string) map[string]string {
	result := make(map[string]string)
	values := r.FindStringSubmatch(s)
	keys := r.SubexpNames()
	for i, value := range values {
		if i != 0 {
			result[keys[i]] = value
		}
	}
	return result

}

func main() {
	initialRecord, err := ParseInitialRecord(sample)
	if err != nil {
		log.Fatalf("%v", err)
	}
	fmt.Printf("%v\n", initialRecord)
}
