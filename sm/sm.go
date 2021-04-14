package sm

type Command struct {
	Index int
	Cmd   string
}

type StatMachine interface {
	Execute(command Command) error
	Retrive(key string) (string, error)
}
