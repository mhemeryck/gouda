package main

import (
	"fmt"
	"log"
)

const sample = `10123                                     0000000550584847241114                                                             000`

func main() {
	record, err := ParseOldBalanceRecord(sample)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	if record.SerialNumber != nil {
		fmt.Printf("%d\n", *record.SerialNumber)
	}
}
