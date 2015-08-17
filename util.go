package main

import (
    "strconv"
    "strings"
)

// parseTincStat creates an object out of tinc log lines
func parseTincStat(loglines []string) (*TincStat, error) {
    // Handle Total Bytes
    bytes_in_line := strings.Fields(loglines[1])
    bytes_out_line := strings.Fields(loglines[2])

    totalbytesin, _ := strconv.ParseInt(bytes_in_line[len(bytes_in_line)-1], 10, 64)
    totalbytesout, _ := strconv.ParseInt(bytes_out_line[len(bytes_out_line)-1], 10, 64)


    // Find and Parse Connections List
    var connections []TincConnection
    inside_connections_stanza := false
    for _, line := range loglines {
        if strings.Contains(line, "Connections:") {
            inside_connections_stanza = true
            continue
        }
        if strings.Contains(line, "End of connections.") {
            inside_connections_stanza = false
            break
        }
        if inside_connections_stanza {
            fields := strings.Fields(line)
            port, _ := strconv.ParseInt(fields[7], 10, 64)
            tinc_conn := TincConnection {
                Name: fields[3],
                Ip: fields[5],
                Port: port,
            }
            connections = append(connections, tinc_conn)
        }



    }

    ts := &TincStat {
        TotalBytesIn: totalbytesin,
        TotalBytesOut: totalbytesout,
        Connections: connections,
    }
    return ts, nil
}
