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

func NewLogEntry(opt int, key, value string) *LogEntry {
	le := &LogEntry{}

	switch opt {
	case SET:
		le.AcceptedValue = "SET " + key + " " + value
	case GET:
		le.AcceptedValue = "GET" + key
	case DELETE:
		le.AcceptedValue = "DELETE" + key
	}

	return le
}

func (le *LogEntry) String() string {
	return strconv.Itoa(le.MinProposal) + ":" + strconv.Itoa(le.AcceptedProposal) + ":" + le.AcceptedValue + ":" + strconv.FormatBool(le.IsCommited)
}
