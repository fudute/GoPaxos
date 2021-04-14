package main

import (
	"flag"
	"fmt"
	"sync"
	"time"

	"github.com/fudute/GoPaxos/driver"
	"github.com/fudute/GoPaxos/utils"
)

var ip = flag.String("ip", "127.0.0.1", "specify the ip address of server")
var port = flag.Int("port", 8000, "port of server")

func main() {
	flag.Parse()

	c := driver.NewClient(*ip, *port)

	var wg sync.WaitGroup
	count := 1000

	wg.Add(count)

	start := time.Now()

	for i := 0; i < count; i++ {
		go func() {
			key := utils.RandString(10)
			value := utils.RandString(10)

			c.Set(key, value)
			wg.Done()
		}()
	}
	wg.Wait()

	fmt.Println(time.Since(start))
}
