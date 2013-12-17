package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os/exec"
)

type UptimeInfo struct {
	One  float64 `json:"one_minute"`
	Five float64 `json:"five_minutes"`
	Ten  float64 `json:"ten_minutes"`
}

func uptime() ([]byte, error) {
	cmd := exec.Command("uptime")
	return cmd.Output()
}

func uptimeServer(w http.ResponseWriter, req *http.Request) {
	output, err := uptime()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ui, err := parseUptimeInfo(output)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.MarshalIndent(ui, " ", "")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(data))
}
