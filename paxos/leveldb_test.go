package paxos

import "testing"

func Test_extractCommand(t *testing.T) {
	value := "17:8:SET ntup xoyhdzdnxe:true"

	command := extractCommand(value)

	if command != "SET ntup xoyhdzdnxe" {
		t.Error("123")
	}
}
