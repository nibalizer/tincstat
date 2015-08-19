package main

import (
	"log"
	"net/http"
)

var (
	outputFile string
)

func main() {
	http.HandleFunc("/tincstat", tincStatServer)
	err := http.ListenAndServe("127.0.0.1:9000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: " + err.Error())
	}
}
