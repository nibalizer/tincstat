package main

import (
	"strconv"
	"strings"
)

// parseUptimeInfo converts the raw uptime command output to an UptimeInfo
// object. It returns an error if any.
func parseUptimeInfo(b []byte) (*UptimeInfo, error) {
	s := string(b)
	// replace commas with spaces, then convert to fields
	uptimes := strings.Fields(strings.Replace(s, ",", " ", -1))

	one, _ := strconv.ParseFloat(uptimes[len(uptimes)-3], 32)
	five, _ := strconv.ParseFloat(uptimes[len(uptimes)-2], 32)
	fifteen, _ := strconv.ParseFloat(uptimes[len(uptimes)-1], 32)
	ui := &UptimeInfo{
		One:     one,
		Five:    five,
		Fifteen: fifteen,
	}
	return ui, nil
}
