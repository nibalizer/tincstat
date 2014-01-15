package main

import "testing"
import "fmt"

func TestParseUptimeInfo(t *testing.T) {

	mymac_string := []byte("19:29  up 30 mins, 4 users, load averages: 0.87 1.03 0.97")
	mymac, err := parseUptimeInfo(mymac_string)
	if err != nil {
		t.Errorf("mac format did not parse: %s parsed as %v\n", mymac_string, mymac)
	}
	if fmt.Sprintf("%0.2f", mymac.One) != "0.87" {
		t.Errorf("One should be 0.87 but was %f\n", mymac.One)
	}
	if fmt.Sprintf("%0.2f", mymac.Five) != "1.03" {
		t.Errorf("One should be 1.03 but was %f\n", mymac.Five)
	}
	if fmt.Sprintf("%0.2f", mymac.Fifteen) != "0.97" {
		t.Errorf("One should be 0.97 but was %f\n", mymac.Fifteen)
	}

	linux_string := []byte("03:58:13 up 155 days, 10:11,  0 users,  load average: 1.69, 1.80, 1.93")
	linux, err := parseUptimeInfo(linux_string)
	if err != nil {
		t.Errorf("linux format did not parse: %s parsed as %v\n", linux_string, linux)
	}
	if fmt.Sprintf("%0.2f", linux.One) != "1.69" {
		t.Errorf("One should be 1.69 but was %f\n", linux.One)
	}
	if fmt.Sprintf("%0.2f", linux.Five) != "1.80" {
		t.Errorf("One should be 1.80 but was %f\n", linux.Five)
	}
	if fmt.Sprintf("%0.2f", linux.Fifteen) != "1.93" {
		t.Errorf("One should be 1.93 but was %f\n", linux.Fifteen)
	}

}
