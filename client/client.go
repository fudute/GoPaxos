package main

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/fudute/GoPaxos/driver"
	"github.com/fudute/GoPaxos/utils"
)

var c *driver.Client

func init() {
	c = driver.NewClient("127.0.0.1", 8000)
}

func TestClient(t *testing.T) {
	key := "name"
	value := "yangxingtai"
	c.Set(key, value)

	name, err := c.Get(key)
	if err != nil {
		t.Error(err)
	}
	if name != value {
		t.Error("test error")
	}
}

func main() {
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
