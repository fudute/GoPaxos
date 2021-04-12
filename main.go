package main

// The following implements the main Go
// package starting up the paxos server

import (
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/fudute/GoPaxos/handlers"
	"github.com/fudute/GoPaxos/paxos"
	log "github.com/sirupsen/logrus"
)

const (
	// HTTP_PORT defines the port value
	// for the Paxos Server service
	HTTP_PORT = "8080"
	RPC_PORT  = ":1234"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {

	//init DB
	paxos.InitDB()
	defer paxos.DB.Close()

	log.WithFields(log.Fields{
		"port": HTTP_PORT,
	}).Info("starting paxos server")
	var wg sync.WaitGroup
	wg.Add(2)

	paxos.InitAcceptor()
	// acceptor serve RPC
	go func() {
		l, e := net.Listen("tcp", RPC_PORT)
		if e != nil {
			log.Fatal("listen error:", e)
		}
		http.Serve(l, nil)
		wg.Done()
	}()

	//connect to acceptor
	paxos.InitProposerNetwork()
	//handle request
	paxos.ProposerHandleRequst()
	//now restful api still not accessiable
	//send a nop request to itself to catch others
	catchUpOthers()

	// RESTful API accessiable
	r := handlers.Router()
	go func() {
		http.ListenAndServe(":"+HTTP_PORT, r)
		wg.Done()
	}()

	wg.Wait()

	// go func() {
	// 	time.Sleep(time.Millisecond * time.Duration(rand.Int()%3000))

	// 	for {
	// 		time.Sleep(time.Second * 3)

	// 		req := paxos.Request{
	// 			Oper: paxos.NOP,
	// 			Done: make(chan error),
	// 		}
	// 		paxos.GetProposerInstance().In <- req

	// 		err := <-req.Done
	// 		if err != nil {
	// 			log.Printf("NOP error: ", err)
	// 		}
	// 	}
	// }()

}

func catchUpOthers() {
	req := paxos.Request{
		Oper: paxos.NOP,
		Done: make(chan error),
	}
	paxos.GetProposerInstance().In <- req

	err := <-req.Done
	if err != nil {
		log.Printf("catch up others NOP error: %v\n", err)
	}
}
