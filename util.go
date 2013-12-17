package main

import (
	"bytes"
	"encoding/csv"
	"strconv"
	"strings"
)

// parseUptimeInfo converts the raw uptime command output to an UptimeInfo
// object. It returns an error if any.
func parseUptimeInfo(b []byte) (*UptimeInfo, error) {
	csvReader := csv.NewReader(bytes.NewReader(b))
	record, err := csvReader.Read()
	if err != nil {
		return nil, err
	}
	uptimes := strings.Split(strings.Trim(record[3], " "), " ")
	one, _ := strconv.ParseFloat(uptimes[2], 32)
	five, _ := strconv.ParseFloat(uptimes[3], 32)
	ten, _ := strconv.ParseFloat(uptimes[4], 32)
	ui := &UptimeInfo{
		One:  one,
		Five: five,
		Ten:  ten,
	}
	return ui, nil
}
