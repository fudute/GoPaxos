package paxos

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/syndtr/goleveldb/leveldb/comparer"
)

const SHIFT = 3

// GenerateProposalNum Generates a proposal number out of minProposalNum and Server ID
func GenerateProposalNum(minProposalNum, ID int) int {
	return (minProposalNum>>SHIFT+1)<<SHIFT + ID
}

// GetPeerNameList Obtains Peer List
// From Environment Variable
func GetPeerNameList() []string {
	return strings.Split(os.Getenv("PEERS"), ",")
}

func GetServerID() int {
	str := os.Getenv("ME")
	id, err := strconv.Atoi(str)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("ServerID = ", id)
	return id
}

// GetNetwork Obtains Network
// From Environment Variable
func GetNetwork() string {
	return os.Getenv("NETWORK") + ":1234"
}

type MyComparator struct {
}

func (cmp MyComparator) Name() string {
	return "my diy cmp"
}
func (cmp MyComparator) Compare(a, b []byte) int {
	i, err1 := strconv.Atoi(string(a))
	j, err2 := strconv.Atoi(string(b))
	if err1 != nil || err2 != nil {
		return comparer.DefaultComparer.Compare(a, b)
	}
	if i == j {
		return 0
	} else if i < j {
		return -1
	}
	return 1
}
func (cmp MyComparator) Separator(dst, a, b []byte) []byte {
	return nil
}

func (cmp MyComparator) Successor(dst, b []byte) []byte {
	return nil
}
func openNewEmptyFile(filePath string) *os.File {
	var file *os.File
	_, err := os.Stat(filePath)
	if err != nil {
		file, err = os.Create(filePath)
		if err != nil {
			log.Fatal("cant create file\n")
		}
	} else {
		file, err = os.OpenFile(filePath, os.O_RDWR|os.O_TRUNC, 0666)
		if err != nil {
			log.Fatal("cant trunc file\n")
		}
	}
	return file
}
