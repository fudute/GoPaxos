package sm

import (
	"reflect"
	"testing"
)

func TestGetKVStatMachine(t *testing.T) {
	tests := []struct {
		name string
		want StatMachine
	}{
		{name: "测试单例模式", want: GetKVStatMachineInstance()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetKVStatMachineInstance(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetKVStatMachine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_kvStatMachine_Execute(t *testing.T) {
	sm := GetKVStatMachineInstance()
	command := "SET name fudute"

	result, err := sm.Execute(command)
	if err != nil {
		t.Errorf(err.Error())
	}
	if result != "" {
		t.Error("result error")
	}

	result, err = sm.Execute("GET name")
	if err != nil {
		t.Errorf(err.Error())
	}
	if result != "fudute" {
		t.Error("result error")
	}
}
