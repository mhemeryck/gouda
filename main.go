package main

import (
	"fmt"
)

const sample = `10000                                     0000000550584847241114                                                             000`

func main() {
	m := groupMap(oldBalanceRecordRegex, sample)
	trimValues(m)
	fmt.Printf("%v", m)
}
