package paxos

import (
	"errors"
	"fmt"
	"log"
	"net/rpc"
	"time"
)

type Proposer struct {
	ServerID  int
	LogIndex  int // 记录最小的没有Chosen的logIndex
	Peers     []*rpc.Client
	peerNames []string
	peerCnt   int

	// 添加一些属性来支持leader
}

var proposer = &Proposer{
	ServerID: GetServerID(),
}

type Request struct {
	Oper  int
	Key   string
	Value string
	Done  chan error
}

type PrepareRequest struct {
	Index       int
	ProposalNum int
}

type PrepareResponse struct {
	AcceptedProposal int
	AcceptedValue    string
}

type AcceptRequest struct {
	Index         int
	ProposalNum   int
	ProposalValue string
}
type AcceptResponse struct {
	MinProposal int
}

type LearnRequest struct {
	Index            int
	AcceptedValue    string
	AcceptedProposal int
}

type LearnResponse struct {
}

const commandSpliter = ";"

func init() {
	proposer.peerNames = GetPeerNameList()
	proposer.peerCnt = len(proposer.peerNames)
	proposer.Peers = make([]*rpc.Client, proposer.peerCnt)
}

func ProposerHandleRequst() {
	go func() {
		for {
			batchReqs := <-GetBatcherInstance().Out
			err := StartNewInstance(batchReqs.Reqs...)
			if err != nil {
				log.Println("Instance error", err)
			}
			// 这里可以选择往done中传不同的参数表示不同的结果
			for i := 0; i < len(batchReqs.Reqs); i++ {
				batchReqs.Reqs[i].Done <- err
			}
			batchReqs.Done <- struct{}{}
		}
	}()
}

func GetProposerInstance() *Proposer {
	return proposer
}

func InitProposerNetwork() {
	// 等待acceptor启动
	time.Sleep(time.Second * 3)
	for index, peerName := range proposer.peerNames {
		addr := fmt.Sprintf("%s:1234", peerName)
		peer, err := rpc.DialHTTP("tcp", addr)
		if err != nil {
			log.Println("Initing Netowork rpc dial error, addr =", addr)
			continue
		}
		proposer.Peers[index] = peer
	}
	for index, peer := range proposer.Peers {
		err := peer.Call("Acceptor.ReDailMe", &proposer.ServerID, &struct{}{})
		if err != nil {
			log.Printf("others cant re-dail to me : %v", err)
		} else {
			log.Printf("node %v re-dail to me succussed ", index)
		}
	}
}

// Prepare starts a Paxos round sending
// a prepare request to all the Paxos
// peers including itself
// 这里发起proposal，直到自己提议的value被chosen，具体的数据传输在doPrepare中完成
// oper取值范围为SET DELETE NOP
func StartNewInstance(reqs ...*Request) error {

	var commands string

	// 这里首先假设key和value只包含26个英文字母，然后使用 ';' 来分割commands中的不同命令
	// 注：在log中，已经使用 ':' 来分割不同的部分，所以在这里不能用 ':'
	for _, req := range reqs {
		var command string
		if req.Oper == SET {
			command = "SET " + req.Key + " " + req.Value
		} else if req.Oper == DELETE {
			command = "DELETE " + req.Key
		} else if req.Oper == NOP {
			command = "NOP"
		} else {
			return ErrorUnkonwCommand
		}
		commands += commandSpliter + command
	}
	commands = commands[1:]

	log.Println("StartNewInstance commands :", commands)

	// 循环获得第一个没有被Chosen的index，直到成功Prepare
	isMeCommited := false
	for !isMeCommited {
		var err error
		le, err := DB.ReadLog(proposer.LogIndex)
		for err == nil && le.IsCommited {
			proposer.LogIndex++
			le, err = DB.ReadLog(proposer.LogIndex)
		}

		if err != nil && err != ErrorNotFound {
			log.Fatal("read error", err)
		}

		isMeCommited, err = DoPrepare(proposer.LogIndex, commands, 0)
		if err != nil {
			return err
		}
		proposer.LogIndex++
	}
	return nil
}

func SendPrepareRequestAndWaitForReply(req *PrepareRequest, done chan struct{}) chan *PrepareResponse {
	out := make(chan *PrepareResponse)

	for _, peer := range proposer.Peers {

		go func(peer *rpc.Client) {
			resp := &PrepareResponse{}
			err := peer.Call("Acceptor.OnPrepare", req, resp)
			if err != nil {
				log.Println("send prepare rpc failed", err)
				return
			}
			select {
			case out <- resp:
			case <-done:
				return
			}
		}(peer)
	}

	return out
}

func SendAcceptRequestAndWaitForReply(req *AcceptRequest, done chan struct{}) chan *AcceptResponse {
	out := make(chan *AcceptResponse)

	for _, peer := range proposer.Peers {

		go func(peer *rpc.Client) {
			resp := &AcceptResponse{}
			err := peer.Call("Acceptor.OnAccept", req, resp)
			if err != nil {
				log.Println("send accept rpc failed", err)
				return
			}
			select {
			case out <- resp:
			case <-done:
				return
			}
		}(peer)
	}

	return out
}

func SendLearnRequest(req *LearnRequest) {
	for _, peer := range proposer.Peers {
		go func(peer *rpc.Client) {
			resp := &LearnResponse{}
			err := peer.Call("Acceptor.OnLearn", req, resp)
			if err != nil {
				log.Println("send learn rpc failed", err)
				return
			}
		}(peer)
	}
}

// DoPrepare可以确定index位置的值
// 这里的value格式为 [SET key value]或者[DELETE key]
// 如果成功提交当前value，返回true，否则返回false
func DoPrepare(index int, value string, minProposal int) (bool, error) {

	log.Printf("[%d] DoPrepare start DoPrepare\n", index)

	proposalNum := GenerateProposalNum(minProposal, proposer.ServerID)

	curValue := value   // 记录当前index的value，有可能之后会变更
	curMaxProposal := 0 // 记录当前看到的最大的accpetedProposal
	preparedPeersCount := 0
	majorityPeersCount := len(proposer.Peers)/2 + 1

	isMeCommited := true

	req := &PrepareRequest{
		Index:       proposer.LogIndex,
		ProposalNum: proposalNum,
	}

	done := make(chan struct{})

	out := SendPrepareRequestAndWaitForReply(req, done)

	for resp := range out {
		preparedPeersCount++
		if resp.AcceptedValue != "" && resp.AcceptedProposal > curMaxProposal {
			curMaxProposal = resp.AcceptedProposal
			curValue = resp.AcceptedValue
			isMeCommited = false
		}
		// Break when majorityPeersCount reached
		if preparedPeersCount >= majorityPeersCount {
			close(done)
			DoAccept(index, proposalNum, curValue)
			break
		}
	}

	if preparedPeersCount < majorityPeersCount {
		close(done)
		return false, errors.New("majority consensus not obtained")
	}

	return isMeCommited, nil
}

// DoAccept starts the accept phase sending
// an accept request to all the Paxos
// peers including itself
func DoAccept(index, proposalNum int, proposalValue string) error {

	log.Printf("[%d] start DoAccept,{proposalNum = %d, value = %s}\n", index, proposalNum, proposalValue)

	acceptedPeersCount := 0
	majorityPeersCount := len(proposer.Peers)/2 + 1

	req := &AcceptRequest{
		Index:         index,
		ProposalNum:   proposalNum,
		ProposalValue: proposalValue,
	}

	done := make(chan struct{})

	out := SendAcceptRequestAndWaitForReply(req, done)

	for resp := range out {
		if resp.MinProposal > proposalNum {
			// 从新prepare，选择更大的
			DoPrepare(index, proposalValue, resp.MinProposal)
			return nil
		}

		acceptedPeersCount++
		// Break when majorityPeersCount reached
		if acceptedPeersCount >= majorityPeersCount {
			// 这里可以直接启动一个协程，向acceptor发送Learn消息
			req := &LearnRequest{
				Index:            index,
				AcceptedProposal: proposalNum,
				AcceptedValue:    proposalValue,
			}
			SendLearnRequest(req)
			break
		}
	}

	if acceptedPeersCount < majorityPeersCount {
		return errors.New("majority consensus not obtained")
	}

	return nil
}
