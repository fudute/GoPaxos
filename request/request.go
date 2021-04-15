package request

import (
	"github.com/fudute/GoPaxos/paxos"
)

const (
	SET = iota
	GET
	DELETE
	NOP
)

type KVRequest struct {
	done chan error
}

type NopReq struct {
	KVRequest
}

type GetReq struct {
	KVRequest
	Key string
}

type SetReq struct {
	KVRequest
	Key   string
	Value string
}

type DelReq struct {
	KVRequest
	Key string
}

func defaultRequest() KVRequest {
	return KVRequest{
		done: make(chan error),
	}
}

func Set(key, value string) paxos.Request {
	return &SetReq{
		KVRequest: defaultRequest(),
		Key:       key,
		Value:     value,
	}
}

func (req *SetReq) GetValue() string {
	return "SET " + req.Key + " " + req.Value
}

func Del(key string) paxos.Request {
	return &DelReq{
		KVRequest: defaultRequest(),
		Key:       key,
	}
}

func (req *DelReq) GetValue() string {
	return "DELETE " + req.Key
}

func Nop() paxos.Request {
	return &NopReq{
		KVRequest: defaultRequest(),
	}
}

func (req *NopReq) GetValue() string {
	return "NOP"
}

func (req *KVRequest) Done() chan error {
	return req.done
}
