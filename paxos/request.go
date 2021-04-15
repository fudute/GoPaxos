package paxos

type Request interface {
	GetValue() string
	Done() chan error
}
