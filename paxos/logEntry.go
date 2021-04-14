package paxos

import "strconv"

const (
	SET = iota
	GET
	DELETE
	NOP
)

type LogEntry struct {
	MinProposal      int
	AcceptedProposal int
	AcceptedValue    string
	IsCommited       bool
}

func (le *LogEntry) String() string {
	return strconv.Itoa(le.MinProposal) + ":" + strconv.Itoa(le.AcceptedProposal) + ":" + le.AcceptedValue + ":" + strconv.FormatBool(le.IsCommited)
}
