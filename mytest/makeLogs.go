package main

import (
	"log"
	"net/http"
)

func main() {
	// logFileName := getCurrentTimeAsFileName()
	logFileName := "log1"
	for _, addr := range serverAddrs {
		resp, err := http.Get(addr + optNop)
		if err != nil || resp.StatusCode != http.StatusOK {
			log.Fatalf("request for syncing logs failed with msg: %v", err)
		}
	}

	for _, addr := range serverAddrs {
		resp, err := http.Get(addr + optPrint + logFileName)
		if err != nil || resp.StatusCode != http.StatusOK {
			log.Fatalf("request for printing logs failed with msg: %v", err)
		}
	}
}
