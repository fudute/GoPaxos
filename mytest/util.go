package main

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"sync"

	// "sync/atomic"
	"testing"
	"time"
)

const (
	smallValLen = 10
	largeValLen = 1000
	serverCnt   = 3
	optSet      = "/store/set/"
	optGet      = "/store/get/"
	optPrint    = "/log/print/"
	optNop      = "/store/nop"
	optIndex    = "/"
	serverIP    = "http://127.0.0.1"
)

var serverAddrs []string = []string{
	serverIP + ":8000",
	serverIP + ":8001",
	serverIP + ":8002",
}

var serverSetUrls []string = []string{
	serverAddrs[0] + optSet,
	serverAddrs[1] + optSet,
	serverAddrs[2] + optSet,
}

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

type KVPair struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

func genRandReader(len int) *bytes.Reader {
	pair := genKVpair(len)
	return genReader(pair)
}
func genReader(pair KVPair) *bytes.Reader {
	bs, err := json.Marshal(pair)
	if err != nil {
		log.Fatalf("cant encode %v to json : %v\n", bs, err)
	}
	reader := bytes.NewReader(bs)
	return reader
}

func randString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 97
		bytes[i] = byte(b)
	}
	return string(bytes)
}

func genKVpair(valuelen int) KVPair {
	if valuelen <= 0 {
		valuelen = 10
	}
	keylen := 5
	return KVPair{
		Key:   randString(keylen),
		Value: randString(valuelen),
	}
}
func genRandReaders(cnt int, len int) []*bytes.Reader {
	readers := make([]*bytes.Reader, cnt)
	for i := 0; i < cnt; i++ {
		readers[i] = genRandReader(len)
	}
	return readers
}

func request(valueLength int, serverUrls []string, b *testing.B) bool {
	repeat := b.N
	readers := genRandReaders(repeat, valueLength)
	serverCnt := len(serverUrls)

	b.ResetTimer()
	for i := 0; i < repeat; i++ {
		http.Post(serverUrls[i%serverCnt], "application/json", readers[i])
		// if err != nil {
		// 	log.Fatalf("async post set value failed with msg: %v \n", err)
		// 	return false
		// }
	}
	return true
}

func asyncRequest(valueLength int, serverUrls []string, b *testing.B) bool {
	repeat := b.N
	// failedCnt := int32(repeat)
	readers := genRandReaders(repeat, valueLength)
	serverCnt := len(serverUrls)

	wg := sync.WaitGroup{}
	wg.Add(repeat)

	b.ResetTimer()
	for i := 0; i < repeat; i++ {
		//time.Sleep(time.Microsecond * 1)
		go func(index int, serverUrl string) {
			defer wg.Done()
			http.Post(serverUrl, "application/json", readers[index])
			// if err != nil {
			// 	log.Fatalf("async post set value faild with msg: %v \n", err)
			// } else {
			// 	atomic.AddInt32(&failedCnt, -1)
			// }
		}(i, serverUrls[i%serverCnt])

	}

	wg.Wait()
	// return failedCnt == 0
	return true
}

// func getCurrentTimeAsFileName() string {
// 	return time.Now().Format("2006-01-02--15-04-05")
// }
