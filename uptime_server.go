package main

import (
    "bufio"
    "encoding/json"
    "errors"
    "fmt"
    "github.com/mitchellh/go-ps"
    "io"
    "log"
    "net/http"
    "os"
    "os/exec"
    "syscall"
)

// UptimeInfo represents the system load average as reported by the uptime command.
type UptimeInfo struct {
	One     float64 `json:"one_minute"`
	Five    float64 `json:"five_minutes"`
	Fifteen float64 `json:"fifteen_minutes"`
}

// findTincPid finds the process of the 'tincd' daemon
func findTincPid() (int, error) {

    procs, err := ps.Processes()
    if err != nil {
        log.Fatalf("findTincPid: %s", err)
    }

    for _, proc := range procs {

        if proc.Executable() == "tincd" {
            //fmt.Println("pid: ", proc.Pid())
            //fmt.Println("ppid: ", proc.PPid())
            //fmt.Println("name: ", proc.Executable())
            //fmt.Println("raw: ", proc)
            return proc.Pid(), err
        }
    }
    return 0, errors.New("findTincPid: Pid not found, is tinc running?")
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
  file, err := os.Open(path)
  if err != nil {
    return nil, err
  }
  defer file.Close()

  var lines []string
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    lines = append(lines, scanner.Text())
  }
  return lines, scanner.Err()
}

// list_contains checks a list for a member
func list_contains(member string, list []string) bool {

    for _,element := range list {
        // index is the index where we are
        // element is the element from someSlice for where we are
        if element == member {
            return true
        }
    }
    return false
}

// uptime executes the uptime command.
func uptime() ([]byte, error) {
	cmd := exec.Command("uptime")
	return cmd.Output()
}

// user12 kills a process with USR1 then USR2
func usr12(pid int) () {
	syscall.Kill(pid, syscall.SIGUSR1)
	syscall.Kill(pid, syscall.SIGUSR2)
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
    // Get tinc pid
    tincPid, err := findTincPid()
    if err != nil {
        log.Fatalf("findTincPid: %s", err)
    }

    // Get first list of lines
    lines1, err := readLines("/var/log/tinc/tinc.log")
    if err != nil {
        log.Fatalf("readLines: %s", err)
    }
    // Send signals
    usr12(tincPid)
    // Confirm flush of data to file
    syscall.Sync()
    // Get second list of lines
    lines2, err := readLines("/var/log/tinc/tinc.log")
    if err != nil {
        log.Fatalf("readLines: %s", err)
    }

    // Print out unique lines in the second set
    for i, line := range lines2 {
        if list_contains(line, lines1) == false {
            fmt.Println(i, line)
        }
    }

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
