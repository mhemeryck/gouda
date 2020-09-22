package main

import (
	"fmt"
	"log"
)

const sample = `0000016032064505                  J. VRANCKEN BV            JVBABE22   00860332194 00005                                       2`

func main() {
	initialRecord, err := ParseInitialRecord(sample)
	if err != nil {
		log.Fatalf("%v", err)
	}
	fmt.Printf("%v\n", initialRecord)
}
