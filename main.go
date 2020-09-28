package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const sample = `10123                                     0000000550584847241114                                                             000`
const filename = "./sample.cod"

func main() {
	// record := OldBalanceRecord{}
	// err := record.Parse(sample)
	// if err != nil {
	// 	log.Fatalf("%v\n", err)
	// }
	// fmt.Printf("%d\n", record.SerialNumber)

	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		log.Fatalf("error opening file %s: %v\n", filename, err)
	}
	scanner := bufio.NewScanner(f)
	records := []Record{}
	var r Record
	for scanner.Scan() {
		line := scanner.Text()
		r, err = Parse(line)
		if err != nil {
			log.Fatalf("error parsing line %s: %v\n", line, err)
		}
		if r != nil {
			records = append(records, r)
		}
	}

	for k, r := range records {
		fmt.Printf("%d: %v\n", k, r)
	}
}
