package paxos

import (
	"fmt"
	"log"
	"net/rpc"
	"sync"

	"github.com/fudute/GoPaxos/sm"
)

type Acceptor struct {
}

var acceptor = &Acceptor{}

func InitAcceptor() {
	rpc.Register(acceptor)
	rpc.HandleHTTP()
}

func (acceptor *Acceptor) ReDailMe(serverID *int, resp *struct{}) error {
	peer, err := rpc.DialHTTP("tcp", fmt.Sprintf("%s:1234", GetPeerNameList()[*serverID]))
	if err != nil {
		log.Println("re-dail rpc dial error, serverID =", *serverID)
		return err
	}
	proposer.Peers[*serverID] = peer
	return nil
}

func GetAcceptorInstance() *Acceptor {
	return acceptor
}

var lo sync.Mutex

// OnPrepare return acceptedProposal and acceptedValue
func (acceptor *Acceptor) OnPrepare(req *PrepareRequest, resp *PrepareResponse) error {

	log.Printf("[%d] receive Prepare {proposalNum = %d}\n", req.Index, req.ProposalNum)

	lo.Lock()
	defer lo.Unlock()

	entry, err := DB.ReadLog(req.Index)
	if err != nil {
		if err == ErrorNotFound {
			entry = &LogEntry{
				AcceptedProposal: 0,
				AcceptedValue:    "",
				MinProposal:      0,
			}
		} else {
			log.Fatal("ReadLog error ", err)
		}
	}
	if req.ProposalNum > entry.MinProposal {
		entry.MinProposal = req.ProposalNum
		// 写回磁盘
		DB.WriteLog(req.Index, entry)
	}
	resp.AcceptedProposal = entry.AcceptedProposal
	resp.AcceptedValue = entry.AcceptedValue

	return nil
}

// OnAccept return minProposal
func (acceptor *Acceptor) OnAccept(req *AcceptRequest, resp *AcceptResponse) error {

	log.Printf("[%d] receive Accept, {proposalNum = %d, proposalValue = %s}\n", req.Index, req.ProposalNum, req.ProposalValue)

	lo.Lock()
	defer lo.Unlock()

	entry, err := DB.ReadLog(req.Index)
	if err != nil && err != ErrorNotFound {
		log.Fatal("ReadLog error ", err)
	}

	if err == ErrorNotFound {
		entry = &LogEntry{}
	}

	if req.ProposalNum >= entry.MinProposal {

		log.Printf("[%d] accepted proposal,{proposalNum = %d, proposalValue = %s}\n", req.Index, req.ProposalNum, req.ProposalValue)

		entry.MinProposal = req.ProposalNum
		entry.AcceptedProposal = req.ProposalNum
		entry.AcceptedValue = req.ProposalValue

		DB.WriteLog(req.Index, entry)
	} else {
		log.Printf("[%d] deny proposal,{proposalNum = %d, minProposal = %d}\n", req.Index, req.ProposalNum, entry.MinProposal)
	}

	resp.MinProposal = entry.MinProposal

	return nil
}

func (acceptor *Acceptor) OnLearn(req *LearnRequest, resp *LearnResponse) error {
	log.Printf("[%d] receive Learn, {LearnValue = %s}\n", req.Index, req.AcceptedValue)

	lo.Lock()
	defer lo.Unlock()

	entry := &LogEntry{}
	entry.AcceptedProposal = req.AcceptedProposal
	entry.MinProposal = req.AcceptedProposal
	entry.AcceptedValue = req.AcceptedValue
	entry.IsCommited = true

	err := DB.WriteLog(req.Index, entry)
	if err != nil {
		log.Fatal("write error: ", err)
	}

	// 在这里执行命令
	sm.GetKVStatMachineInstance().Execute(req.AcceptedValue)

	return err
}
