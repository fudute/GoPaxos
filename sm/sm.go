package sm

type StatMachine interface {
	Execute(command string) (string, error)
}
