package paxos

import (
	"github.com/fudute/GoPaxos/sm"
)

// store used to store & retrieve string entries

type Store interface {
	WriteLog(index int, entry *LogEntry) error
	ReadLog(index int) (*LogEntry, error)
	Restore(p *Proposer, a *Acceptor, sm sm.StatMachine) error // 重启后读取文件恢复状态
	PrintLog(fileName string)
	Close() error
}

var DB Store
