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

// TincStat represents the status of the tinc daemon

type TincStat struct {
    TotalBytesIn int64 `json:"total_bytes_in"`
    TotalBytesOut int64 `json:"total_bytes_out"`
    Connections []TincConnection `json:"connections"`
}

type TincConnection struct {
    Name string `json:"name"`
    Ip string `json:"ip"`
    Port int64 `json:"port"`
}



// findTincPid finds the process of the 'tincd' daemon
// TODO: reading the pidfile might be smarter
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


// tincStatServer serves data pulled from the tinc log file
// Tinc logs connection and network information after getting a USR1 and USR2
// The following output is current:
//
//{
//  "total_bytes_in": 115324,
//  "total_bytes_out": 67990,
//  "connections": [
//    {
//      "name": "some_random_node",
//      "ip": "192.0.2.15",
//      "port": 2003
//    }
//  ]
//}
func tincStatServer(w http.ResponseWriter, req *http.Request) {
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

    // Print out and save unique lines in the second set
    var loglines []string
    for i, line := range lines2 {
        if list_contains(line, lines1) == false {
            fmt.Println(i, line)
            loglines = append(loglines, line)
        }
    }

    // Convert the raw loglines output to a tincstat object
    ts, err := parseTincStat(loglines)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    // Create the JSON representation of tinc status
    data, err := json.MarshalIndent(ts, " ", "")
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    // Write the HTTP response headers and body.
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    io.WriteString(w, string(data))
}
