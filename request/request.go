package request

import "log"

const (
	SET = iota
	GET
	DELETE
	NOP
)

type KVRequest struct {
	Oper  int
	Key   string
	Value string
	done  chan error
}

func (req *KVRequest) GetValue() string {

	var command string
	if req.Oper == SET {
		command = "SET " + req.Key + " " + req.Value
	} else if req.Oper == DELETE {
		command = "DELETE " + req.Key
	} else if req.Oper == NOP {
		command = "NOP"
	} else {
		// TODO 这里需要再斟酌一下
		log.Println("Unknow command")
		return ""
	}
	return command
}

func (req *KVRequest) Done() chan error {
	return req.done
}
