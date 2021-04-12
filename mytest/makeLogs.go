package main

import (
	"fmt"
	"net/http"
)

func main() {
	// logFileName := getCurrentTimeAsFileName()
	logFileName := "log1"
	for _, addr := range serverAddrs {
		resp, err := http.Get(addr + optNop)
		if err != nil || resp.StatusCode != http.StatusOK {
			_ = fmt.Errorf("request for syncing logs failed with msg: %v \n", err)
		}
	}

	for _, addr := range serverAddrs {
		resp, err := http.Get(addr + optPrint + logFileName)
		if err != nil || resp.StatusCode != http.StatusOK {
			_ = fmt.Errorf("request for printing logs failed with msg: %v \n", err)
		}
	}
}
