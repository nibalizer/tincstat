package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os/exec"
)

// UptimeInfo represents the system load average as reported by the uptime command.
type UptimeInfo struct {
	One     float64 `json:"one_minute"`
	Five    float64 `json:"five_minutes"`
	Fifteen float64 `json:"fifteen_minutes"`
}

// uptime executes the uptime command.
func uptime() ([]byte, error) {
	cmd := exec.Command("uptime")
	return cmd.Output()
}

// uptimeServer servers the system load average as reported by the uptime
// command. It returns the system uptime in the following format:
//
//   {
//     "one_minute": 1.0199999809265137,
//     "five_minutes": 1.2100000381469727,
//     "fifteen_minutes": 1.2300000190734863
//   }
//
func uptimeServer(w http.ResponseWriter, req *http.Request) {
	output, err := uptime()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Convert the raw uptime output to an UptimeInfo object.
	ui, err := parseUptimeInfo(output)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create the JSON representation of the system uptime.
	data, err := json.MarshalIndent(ui, " ", "")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write the HTTP response headers and body.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(data))
}
