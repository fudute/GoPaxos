package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

var port *int
var addr *string

var wg sync.WaitGroup

func randString(n int) string {
	lenght := rand.Int()%n + 1

	str := make([]byte, lenght)
	for i := 0; i < lenght; i++ {
		str[i] = byte(rand.Int()%26) + 'a'
	}
	return string(str)
}

func testGet(n int) {
	wg.Add(n)
	start := time.Now()
	for i := 0; i < n; i++ {
		go func() {
			_, err := http.Get(fmt.Sprintf("http://%s:%d/get/name", *addr, *port))
			if err != nil {
				log.Fatal(err)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	fmt.Println(time.Since(start))

}

type KVPair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func testSet(n int) {

	wg.Add(n)
	start := time.Now()
	for i := 0; i < n; i++ {

		kvp := &KVPair{
			Key:   randString(10),
			Value: randString(10),
		}

		bs, err := json.Marshal(kvp)
		if err != nil {
			log.Fatal(err)
		}

		time.Sleep(time.Millisecond * 1)
		go func() {
			_, err := http.Post(fmt.Sprintf("http://%s:%d/store/set", *addr, *port), "application/json", bytes.NewReader(bs))
			if err != nil {
				log.Fatal(err)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	fmt.Println(time.Since(start))
}

// "set key value"

func main() {
	addr = flag.String("addr", "127.0.0.1", "address of server")
	port = flag.Int("port", 8000, "port to server")
	flag.Parse()

	fmt.Println("port = ", *port)
	rand.Seed(time.Now().Unix())
	n := 1000
	testSet(n)
}
