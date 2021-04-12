package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const ()

var HomePath string
var logAbsPaths []string
var relativePaths []string = []string{
	"/paxos_logs/n0/log1",
	"/paxos_logs/n1/log1",
	"/paxos_logs/n2/log1",
}
var nodeCnt int
var fileScaners []*bufio.Scanner
var major int

func init() {
	HomePath = os.Getenv("HOME")
	nodeCnt = len(relativePaths)
	logAbsPaths = make([]string, 0, nodeCnt)
	fileScaners = make([]*bufio.Scanner, 0, nodeCnt)
	major = (nodeCnt)/2 + 1
}

type Log struct {
	Index      string
	IsCommited string
	Command    string
}

func getLogFromString(a string) Log {
	strs := strings.Split(a, " ")
	if len(strs) != 3 {
		_ = fmt.Errorf("line parse error %v", a)
	}
	command := strs[2]
	for _, str := range strs[3:] {
		command += " " + str
	}
	return Log{
		Index:      strs[0],
		IsCommited: strs[1],
		Command:    command,
	}
}
func getMajorValue(logs []Log) string {

	var m map[string]int = make(map[string]int)
	for _, log := range logs {
		m[log.Command]++
		if m[log.Command] >= major {
			return log.Command
		}
	}
	return ""
}
func checkLogConsistency(logStrs []string) (bool, string) {
	logs := make([]Log, 0, nodeCnt)
	for _, str := range logStrs {
		log := getLogFromString(str)
		logs = append(logs, log)
	}
	majorValue := getMajorValue(logs)
	if majorValue == "" || majorValue == "UNKNOWN" {
		return false, majorValue
	}
	for _, log := range logs {
		if log.IsCommited == "true" && log.Command != majorValue {
			return false, majorValue
		}
	}
	return true, majorValue
}

func main() {
	for i := 0; i < nodeCnt; i++ {
		logAbsPaths = append(logAbsPaths, HomePath+relativePaths[i])
		fmt.Println(logAbsPaths[i])
		file, err := os.OpenFile(logAbsPaths[i], os.O_RDONLY, 0622)
		if err != nil {
			_ = fmt.Errorf("cant open file %v\n", err)
		}
		defer file.Close()
		fileScaners = append(fileScaners, bufio.NewScanner(file))
	}

	eofFlags := make([]bool, nodeCnt)
	remainCnt := nodeCnt
	index := 0
	for {
		strs := make([]string, 0, nodeCnt)
		for i := 0; i < nodeCnt; i++ {
			if fileScaners[i].Scan() {
				str := fileScaners[i].Text()
				strs = append(strs, str)
			} else {
				eofFlags[i] = true
				remainCnt--
				fmt.Println("--------------------------over--------------------------")
				fmt.Println("the consistency of the remaining part is not guaranteed!")
				fmt.Println("--------------------------over--------------------------")
				fmt.Printf("node %v eof\n", i)
				break
			}

		}
		if remainCnt != nodeCnt {
			break
		}
		ok, val := checkLogConsistency(strs)
		if !ok {
			var errStr string
			for _, str := range strs {
				errStr += str + "\n"
			}
			_ = fmt.Errorf("index %v doesn't match,value :%v", index, errStr)
		}
		fmt.Printf("consensus for index %v     [%v]\n", index, val)
		index++
	}
	for remainCnt > 0 {
		for i := 0; i < nodeCnt; i++ {
			if eofFlags[i] {
				continue
			}
			if fileScaners[i].Scan() {
				str := fileScaners[i].Text()
				fmt.Printf("node %v :%v\n", i, str)
			} else {
				eofFlags[i] = true
				remainCnt--
				fmt.Printf("node %v eof\n", i)
			}
		}
		index++
	}

}
