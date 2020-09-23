package main

import (
	"testing"
	"time"
)

func TestParseInitialRecord(t *testing.T) {
	sample := `0000013020912605        YjeybrNhwgMichael Campbell          BBRUBEBB   03155032542                                             2`
	r, err := ParseInitialRecord(sample)
	if err != nil {
		t.Fatalf("could not parse sample: %v", err)
	}
	if !r.CreationDate.Equal(time.Date(2009, 2, 13, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("wrong date parsed")
	}
	if r.BankIdentificationNumber != 126 {
		t.Fatalf("incorrect bank identification code parsed")
	}
	if r.Reference != "YjeybrNhwg" {
		t.Fatalf("wrong reference parsed")
	}
	if r.AccountHolderReference != 3155032542 {
		t.Fatalf("issue parsing account holder reference")
	}
	if r.BIC != "BBRUBEBB" {
		t.Fatalf("issue parsing BIC code")
	}
	if r.IsDuplicate == true {
		t.Fatalf("issue parsing duplicate code")
	}
	if r.Addressee != "Michael Campbell" {
		t.Fatalf("issue parsing addressee")
	}
}

func TestParseOldBalanceRecord(t *testing.T) {
	sample := `10000                                     0000000550584847241114                                                             000`
	r, err := ParseOldBalanceRecord(sample)
	if err != nil {
		t.Fatalf("could not parse sample: %v", err)
	}
	t.Logf("%v\n", r)
}
