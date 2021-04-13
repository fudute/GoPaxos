package sm

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

type kvStatMachine struct {
	m map[string]string
}

// 单例模式
var sm StatMachine = &kvStatMachine{
	m: map[string]string{},
}

func GetKVStatMachineInstance() StatMachine {
	return sm
}

func (kvSM *kvStatMachine) Execute(command string) (string, error) {
	if len(command) == 0 {
		fmt.Println("empty command")
		return "", nil
	}
	log.Printf("execute command: %s\n", command)

	tokens := strings.Split(command, " ")

	if tokens[0] == "SET" {
		if len(tokens) < 3 {
			return "", errors.New("command format Error")
		}
		kvSM.m[tokens[1]] = tokens[2]
		return "", nil

	} else if tokens[0] == "GET" {
		if len(tokens) < 2 {
			return "", errors.New("command format Error")
		}
		return kvSM.m[tokens[1]], nil

	} else if tokens[0] == "DELETE" {
		if len(tokens) < 2 {
			return "", errors.New("command format Error")
		}
		delete(kvSM.m, tokens[1])
		return "", nil
	} else if tokens[0] == "NOP" {
		return "", nil
	}
	return "", errors.New("unkonw command")
}
