package sm

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
)

// 需要保证状态机是按序执行的，所以说需要实现一个滑动窗口机制

type kvStatMachine struct {
	m         map[string]string
	nextIndex int
	cmds      map[int]Command
	mu        sync.Mutex
}

// 单例模式
var sm StatMachine = &kvStatMachine{
	m:         map[string]string{},
	nextIndex: 0,
	cmds:      make(map[int]Command),
}

func GetKVStatMachineInstance() StatMachine {
	return sm
}

func (kvSM *kvStatMachine) doExecute(command string) error {
	if len(command) == 0 {
		fmt.Println("empty command")
		return nil
	}
	log.Printf("execute command: %s\n", command)

	tokens := strings.Split(command, " ")

	if tokens[0] == "SET" {
		if len(tokens) < 3 {
			return errors.New("command format Error")
		}
		kvSM.m[tokens[1]] = tokens[2]
		return nil

	} else if tokens[0] == "DELETE" {
		if len(tokens) < 2 {
			return errors.New("command format Error")
		}
		delete(kvSM.m, tokens[1])
		return nil
	} else if tokens[0] == "NOP" {
		return nil
	}
	return errors.New("unkonw command")
}
func (kvSM *kvStatMachine) Execute(cmd Command) error {

	kvSM.mu.Lock()
	defer kvSM.mu.Unlock()

	kvSM.cmds[cmd.Index] = cmd

	for {
		cmd, ok := kvSM.cmds[kvSM.nextIndex]
		if ok {
			delete(kvSM.cmds, kvSM.nextIndex)
			err := kvSM.doExecute(cmd.Cmd)
			if err != nil {
				return err
			}
			kvSM.nextIndex++
		} else {
			break
		}
	}
	return nil
}

func (kvSM *kvStatMachine) Retrive(key string) (string, error) {
	val, ok := kvSM.m[key]
	if ok {
		return val, nil
	}
	return "", nil
}
